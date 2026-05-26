def test_root_endpoint(host):
    result = host.run("curl -sf -o /dev/null -w '%{http_code}' http://127.0.0.1/")
    assert result.stdout.strip() == "200"


def test_notes_endpoint_json(host):
    result = host.run(
        "curl -sf -o /dev/null -w '%{http_code}' -H 'Accept: application/json' http://127.0.0.1/notes"
    )
    assert result.stdout.strip() == "200"


def test_notes_endpoint_html(host):
    result = host.run(
        "curl -sf -o /dev/null -w '%{http_code}' -H 'Accept: text/html' http://127.0.0.1/notes"
    )
    assert result.stdout.strip() == "200"


def test_health_blocked_by_nginx(host):
    result = host.run("curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1/health/alive")
    assert result.stdout.strip() == "404"


def test_health_direct(host):
    result = host.run("curl -sf http://127.0.0.1:5000/health/alive")
    assert "OK" in result.stdout
