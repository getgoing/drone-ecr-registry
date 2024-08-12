load('@common//:images.star', 'find_image')
load(
    '@common//:registries.star',
    'find_registry',
    'registry_uri',
)
load(
    '@common//:runners.star',
    'find_runner',
    'runner_env',
    'runner_node',
)
load('@common//:steps.star', 'notify_author')
load('@common//:utils.star', 'retrieve_parameter')

REGISTRY = find_registry('infrastructure')

INFRA_RUNNER = find_runner(environment='infrastructure', type='kubernetes', rclass='spot')

def main(ctx):
    return [
        retrieve_parameter('DRONE_SLACK_BOT_TOKEN'),
        retrieve_parameter('DRONE_PEOPLEFORCE_API_KEY'),
        pr_pipeline(ctx),
        build_pipeline(ctx),
    ]

def pr_pipeline(ctx):
    return {
        'kind': 'pipeline',
        'name': 'pull request',
        'steps': [
            generate_tags_file(ctx),
            build_and_push(dry_run=True),
            notify_author(
                {'from_secret': 'drone_slack_bot_token'},
                {'from_secret': 'drone_peopleforce_api_key'}
            ),
        ],
        'trigger': {
            'event': ['pull_request'],
        },
        'node': runner_node(INFRA_RUNNER),
    }

def build_pipeline(ctx):
    return {
        'kind': 'pipeline',
        'name': 'build and push',
        'steps': [
            generate_tags_file(ctx),
            build_and_push(),
            notify_author(
                {'from_secret': 'drone_slack_bot_token'},
                {'from_secret': 'drone_peopleforce_api_key'}
            ),
        ],
        'trigger': {
            'event': ['push'],
            'branch': ['master'],
        },
        'node': runner_node(INFRA_RUNNER),
    }

def build_and_push(dry_run=False):
    if dry_run:
        name = 'build docker image (dry-run)'
    else:
        name = 'build and push docker image'
        dry_run = False

    return {
        'name': name,
        'image': find_image('plugins/ecr'),
        'settings': {
            'registry': registry_uri(REGISTRY),
            'repo': 'drone-plugin/ecr-registry',
            'custom_dns': '169.254.169.253',
            'dry_run': dry_run,
            'no_cache': True,
        },
        'environment': runner_env(INFRA_RUNNER),
    }

def generate_tags_file(ctx):
    commit_sha = ctx.build.commit[:6]

    return {
        'name': 'generate tags file',
        'image': find_image('alpine'),
        'commands': [
            'echo -n "$(cat version),$DRONE_BUILD_NUMBER,latest,{}" > .tags'.format(commit_sha),
        ],
    }
