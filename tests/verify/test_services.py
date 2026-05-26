def test_nginx_running(host):
    assert host.service("nginx").is_running
    assert host.service("nginx").is_enabled


def test_mywebapp_running(host):
    assert host.service("mywebapp").is_running
    assert host.service("mywebapp").is_enabled


def test_postgresql_running(host):
    assert host.service("postgresql").is_running
