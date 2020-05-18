package crossplane_test

/*
def test_build_nested_and_multiple_args():
    payload = [
        {
            "directive": "events",
            "args": [],
            "block": [
                {
                    "directive": "worker_connections",
                    "args": ["1024"]
                }
            ]
        },
        {
            "directive": "http",
            "args": [],
            "block": [
                {
                    "directive": "server",
                    "args": [],
                    "block": [
                        {
                            "directive": "listen",
                            "args": ["127.0.0.1:8080"]
                        },
                        {
                            "directive": "server_name",
                            "args": ["default_server"]
                        },
                        {
                            "directive": "location",
                            "args": ["/"],
                            "block": [
                                {
                                    "directive": "return",
                                    "args": ["200", "foo bar baz"]
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    ]
    built = crossplane.build(payload, indent=4, tabs=False)
    assert built == '\n'.join([
        'events {',
        '    worker_connections 1024;',
        '}',
        'http {',
        '    server {',
        '        listen 127.0.0.1:8080;',
        '        server_name default_server;',
        '        location / {',
        "            return 200 'foo bar baz';",
        '        }',
        '    }',
        '}'
    ])


def test_build_with_comments():
    payload = [
        {
            "directive": "events",
            "line": 1,
            "args": [],
            "block": [
                {
                    "directive": "worker_connections",
                    "line": 2,
                    "args": ["1024"]
                }
            ]
        },
        {
            "directive": "#",
            "line": 4,
            "args": [],
            "comment": "comment"
        },
        {
            "directive": "http",
            "line": 5,
            "args": [],
            "block": [
                {
                    "directive": "server",
                    "line": 6,
                    "args": [],
                    "block": [
                        {
                            "directive": "listen",
                            "line": 7,
                            "args": ["127.0.0.1:8080"]
                        },
                        {
                            "directive": "#",
                            "line": 7,
                            "args": [],
                            "comment": "listen"
                        },
                        {
                            "directive": "server_name",
                            "line": 8,
                            "args": ["default_server"]
                        },
                        {
                            "directive": "location",
                            "line": 9,
                            "args": ["/"],
                            "block": [
                                {
                                    "directive": "#",
                                    "line": 9,
                                    "args": [],
                                    "comment": "# this is brace"
                                },
                                {
                                    "directive": "#",
                                    "line": 10,
                                    "args": [],
                                    "comment": " location /"
                                },
                                {
                                    "directive": "#",
                                    "line": 11,
                                    "args": [],
                                    "comment": " is here"
                                },
                                {
                                    "directive": "return",
                                    "line": 12,
                                    "args": ["200", "foo bar baz"]
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    ]
    built = crossplane.build(payload, indent=4, tabs=False)
    assert built == '\n'.join([
        'events {',
        '    worker_connections 1024;',
        '}',
        '#comment',
        'http {',
        '    server {',
        '        listen 127.0.0.1:8080; #listen',
        '        server_name default_server;',
        '        location / { ## this is brace',
        '            # location /',
        '            # is here',
        "            return 200 'foo bar baz';",
        '        }',
        '    }',
        '}'
    ])


def test_build_starts_with_comments():
    payload = [
        {
            "directive": "#",
            "line": 1,
            "args": [],
            "comment": " foo"
        },
        {
            "directive": "user",
            "line": 5,
            "args": ["root"]
        }
    ]
    built = crossplane.build(payload, indent=4, tabs=False)
    assert built == '# foo\nuser root;'


def test_build_with_quoted_unicode():
    payload = [
        {
            "directive": "env",
            "line": 1,
            "args": ["русский текст"],
        }
    ]
    built = crossplane.build(payload, indent=4, tabs=False)
    assert built == u"env 'русский текст';"


def test_build_multiple_comments_on_one_line():
    payload = [
        {
            "directive": "#",
            "line": 1,
            "args": [],
            "comment": "comment1"
        },
        {
            "directive": "user",
            "line": 2,
            "args": ["root"]
        },
        {
            "directive": "#",
            "line": 2,
            "args": [],
            "comment": "comment2"
        },
        {
            "directive": "#",
            "line": 2,
            "args": [],
            "comment": "comment3"
        }
    ]
    built = crossplane.build(payload, indent=4, tabs=False)
    assert built == '#comment1\nuser root; #comment2 #comment3'



def test_build_files_with_missing_status_and_errors(tmpdir):
    assert len(tmpdir.listdir()) == 0
    payload = {
        "config": [
            {
                "file": "nginx.conf",
                "parsed": [
                    {
                        "directive": "user",
                        "line": 1,
                        "args": ["nginx"],
                    }
                ]
            }
        ]
    }
    crossplane.builder.build_files(payload, dirname=tmpdir.strpath)
    built_files = tmpdir.listdir()
    assert len(built_files) == 1
    assert built_files[0].strpath == os.path.join(tmpdir.strpath, 'nginx.conf')
    assert built_files[0].read_text('utf-8') == 'user nginx;\n'


def test_build_files_with_unicode(tmpdir):
    assert len(tmpdir.listdir()) == 0
    payload = {
        "status": "ok",
        "errors": [],
        "config": [
            {
                "file": "nginx.conf",
                "status": "ok",
                "errors": [],
                "parsed": [
                    {
                        "directive": "user",
                        "line": 1,
                        "args": [u"測試"],
                    }
                ]
            }
        ]
    }
    crossplane.builder.build_files(payload, dirname=tmpdir.strpath)
    built_files = tmpdir.listdir()
    assert len(built_files) == 1
    assert built_files[0].strpath == os.path.join(tmpdir.strpath, 'nginx.conf')
    assert built_files[0].read_text('utf-8') == u'user 測試;\n'


def test_compare_parsed_and_built_simple(tmpdir):
    compare_parsed_and_built('simple', 'nginx.conf', tmpdir)


def test_compare_parsed_and_built_messy(tmpdir):
    compare_parsed_and_built('messy', 'nginx.conf', tmpdir)


def test_compare_parsed_and_built_messy_with_comments(tmpdir):
    compare_parsed_and_built('with-comments', 'nginx.conf', tmpdir, comments=True)


def test_compare_parsed_and_built_empty_map_values(tmpdir):
    compare_parsed_and_built('empty-value-map', 'nginx.conf', tmpdir)


def test_compare_parsed_and_built_russian_text(tmpdir):
    compare_parsed_and_built('russian-text', 'nginx.conf', tmpdir)


def test_compare_parsed_and_built_quoted_right_brace(tmpdir):
    compare_parsed_and_built('quoted-right-brace', 'nginx.conf', tmpdir)


def test_compare_parsed_and_built_directive_with_space(tmpdir):
    compare_parsed_and_built('directive-with-space', 'nginx.conf', tmpdir)
*/


/*
import os

from crossplane.compat import basestring
from crossplane.parser import parse
from crossplane.builder import build, _enquote

here = os.path.dirname(__file__)


def assert_equal_payloads(a, b, ignore_keys=()):
    assert type(a) == type(b)
    if isinstance(a, list):
        assert len(a) == len(b)
        for args in zip(a, b):
            assert_equal_payloads(*args, ignore_keys=ignore_keys)
    elif isinstance(a, dict):
        keys = set(a.keys()) | set(b.keys())
        keys.difference_update(ignore_keys)
        for key in keys:
            assert_equal_payloads(a[key], b[key], ignore_keys=ignore_keys)
    elif isinstance(a, basestring):
        assert _enquote(a) == _enquote(b)
    else:
        assert a == b


def compare_parsed_and_built(conf_dirname, conf_basename, tmpdir, **kwargs):
    original_dirname = os.path.join(here, 'configs', conf_dirname)
    original_path = os.path.join(original_dirname, conf_basename)
    original_payload = parse(original_path, **kwargs)
    original_parsed = original_payload['config'][0]['parsed']

    build1_config = build(original_parsed)
    build1_file = tmpdir.join('build1.conf')
    build1_file.write_text(build1_config, encoding='utf-8')
    build1_payload = parse(build1_file.strpath, **kwargs)
    build1_parsed = build1_payload['config'][0]['parsed']

    assert_equal_payloads(original_parsed, build1_parsed, ignore_keys=['line'])

    build2_config = build(build1_parsed)
    build2_file = tmpdir.join('build2.conf')
    build2_file.write_text(build2_config, encoding='utf-8')
    build2_payload = parse(build2_file.strpath, **kwargs)
    build2_parsed = build2_payload['config'][0]['parsed']

    assert build1_config == build2_config
    assert_equal_payloads(build1_parsed, build2_parsed, ignore_keys=[])
*/