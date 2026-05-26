def test_nginx_site_enabled(host):
    assert host.file("/etc/nginx/sites-enabled/mywebapp").exists


def test_nginx_config_valid(host):
    cmd = host.run("nginx -t")
    assert cmd.rc == 0


def test_nginx_access_log(host):
    assert host.file("/var/log/nginx/mywebapp_access.log").exists
