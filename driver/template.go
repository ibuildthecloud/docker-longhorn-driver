package driver

const (
	DockerComposeTemplate = `
replica:
    scale: 2
    image: cjellick/longhorn:dev
    command:
    - launch
    - replica
    - --listen
    - 0.0.0.0:9502
    - --sync-agent=false
    - /var/lib/rancher/longhorn/$VOLUME_NAME
    volumes:
    - /var/lib/rancher/longhorn/$VOLUME_NAME
    labels:
        io.rancher.sidekicks: replica-agent, sync-agent
        io.rancher.container.hostname_override: container_name
        io.rancher.scheduler.affinity:container_label_ne: io.rancher.stack_service.name=$${stack_name}/$${service_name}
        io.rancher.scheduler.affinity:container_soft: $LONGHORN_DRIVER_CONTAINER
    metadata:
        longhorn:
            volume_name: $VOLUME_NAME
            volume_size: $VOLUME_SIZE
    health_check:
        healthy_threshold: 1
        unhealthy_threshold: 3
        interval: 5000
        port: 8199
        request_line: GET /replica/status HTTP/1.0
        response_timeout: 50000
        strategy: recreateOnQuorum
        recreate_on_quorum_strategy_config:
            quorum: 1

sync-agent:
    image: cjellick/longhorn:dev
    net: container:replica
    working_dir: /var/lib/rancher/longhorn/$VOLUME_NAME
    volumes_from:
    - replica
    command:
    - longhorn
    - sync-agent
    - --listen
    - 0.0.0.0:9504

replica-agent:
    image: cjellick/longhorn:dev
    net: container:replica
    metadata:
        longhorn:
            volume_name: $VOLUME_NAME
            volume_size: $VOLUME_SIZE
    command:
    - longhorn-agent
    - --replica

controller:
    image: cjellick/longhorn:dev
    command:
    - launch
    - controller
    - --listen
    - 0.0.0.0:9501
    - --frontend
    - tcmu
    - $VOLUME_NAME
    privileged: true
    volumes:
    - /dev:/host/dev
    - /lib/modules:/lib/modules:ro
    labels:
        io.rancher.sidekicks: controller-agent
        io.rancher.container.hostname_override: container_name
        io.rancher.scheduler.affinity:container: $LONGHORN_DRIVER_CONTAINER
    metadata:
        longhorn:
          volume_name: $VOLUME_NAME
    health_check:
        healthy_threshold: 1
        interval: 5000
        port: 8199
        request_line: GET /controller/status HTTP/1.0
        response_timeout: 5000
        strategy: none
        unhealthy_threshold: 2

controller-agent:
    image: cjellick/longhorn:dev
    net: container:controller
    metadata:
        longhorn:
          volume_name: $VOLUME_NAME
    command:
    - longhorn-agent
    - --controller
`
)
