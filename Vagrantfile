Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-24.04"

  config.vm.network "forwarded_port", guest: 80, host: 8080

  config.vm.provider "virtualbox" do |vb|
    vb.memory = 1024
    vb.cpus = 1
  end

  config.vm.provider "vmware_desktop" do |vmw|
    vmw.memory = 1024
    vmw.cpus = 1
  end

  config.vm.provision "shell", inline: <<-SHELL
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
    bash /vagrant/deploy/setup_target.sh
  SHELL
end
