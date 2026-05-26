Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-24.04"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = 1024
    vb.cpus = 1
  end

  config.vm.provider "vmware_desktop" do |vmw|
    vmw.memory = 1024
    vmw.cpus = 1
  end

  config.vm.define "runner" do |runner|
    runner.vm.hostname = "runner"
    runner.vm.network "private_network", ip: "192.168.56.11"

    runner.vm.provision "shell", inline: <<-SHELL
      set -euo pipefail
      bash /vagrant/deploy/setup_runner.sh
    SHELL
  end

  config.vm.define "target" do |target|
    target.vm.hostname = "target"
    target.vm.network "private_network", ip: "192.168.56.10"
    target.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"

    target.vm.provision "shell", inline: <<-SHELL
      set -euo pipefail
      if [[ -f /vagrant/.env ]]; then
        set -a
        source /vagrant/.env
        set +a
      fi
      export IMAGE_REF="${IMAGE_REF:-mywebapp:local}"
      export DB_USER="${DB_USER:-mywebapp}"
      export DB_PASSWORD="${DB_PASSWORD:-change-me}"
      export DB_NAME="${DB_NAME:-mywebapp}"
      export STUDENT_PASSWORD="${STUDENT_PASSWORD:-change-me}"
      mkdir -p /tmp/mywebapp-deploy
      cp /vagrant/deploy/config.yaml.tmpl /tmp/mywebapp-deploy/
      cp /vagrant/deploy/mywebapp.service.tmpl /tmp/mywebapp-deploy/
      cp /vagrant/deploy/nginx.conf /tmp/mywebapp-deploy/
      cp /vagrant/deploy/sudoers-operator /tmp/mywebapp-deploy/
      cp /vagrant/deploy/sudoers-student /tmp/mywebapp-deploy/
      bash /vagrant/deploy/setup_target.sh
      if [[ -f /vagrant/.runner_pubkey ]]; then
        mkdir -p /home/student/.ssh
        chmod 700 /home/student/.ssh
        grep -qF "$(cat /vagrant/.runner_pubkey)" /home/student/.ssh/authorized_keys 2>/dev/null || \
          cat /vagrant/.runner_pubkey >> /home/student/.ssh/authorized_keys
        chmod 600 /home/student/.ssh/authorized_keys
        chown -R student:student /home/student/.ssh
      fi
      docker build -t mywebapp:local /vagrant
      systemctl restart mywebapp.service
    SHELL
  end
end
