def test_postgres_listens_localhost(host):
    assert host.socket("tcp://127.0.0.1:5432").is_listening


def test_postgres_not_exposed(host):
    assert not host.socket("tcp://0.0.0.0:5432").is_listening
