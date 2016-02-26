# snap collector plugin - dbi

This plugin connects to various databases, executes SQL statements and reads back the results. 

The plugin is a generic plugin. You can configure how each column is to be interpreted and the plugin will generate one or more data sets from each row returned according. These rules are defined in separated json file (read more about this in [Configuration and Usage](#configuration-and-usage) section).

Depending on the configuration, the returned values are converted into metrics.

The plugin is used in the [snap framework] (http://github.com/intelsdi-x/snap).				

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Setfile fields](#setfile-fields)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

- Linux system
- Access to database (currently the following SQL Drivers are supported: **MySQL**, **PostgreSQL**)

### Installation

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-dbi  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-dbi.git
```

Build the snap dbi plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

* Create configuration file (called as a setfile) in which will be defined databases, queries and rules how interpret the results, see exemplary in [examples/configs/setfiles](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/master/examples/configs/setfiles)

* Set up field `setfile` in Global Config as a path to dbi plugin configuration file, see exemplary Snap Global Config: in [examples/configs/snap-config-sample.json] (https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/master/examples/configs/snap-config-sample.json)
 
Notice that this plugin is a generic plugin, i.e. it cannot work without configuration, because there is no reasonable default behavior.

## Documentation

### Setfile fields

* **queries** - contains all defined queries put in query block which includes:
	*  **name** - identify query block, needs to be unique
	*  **statement** - SQL statement to be executed
	*  **results** - block which defines results of statement
* **results** - contains how the returned data should be interpreted, including:
	 * **name** - name of result, acceptable empty if only one result is defined; in other case must be given in order to distinguish results
	* **instance_from** - name of column whose values will be used to specify an instance
	* **instance_prefix** - prepended prefix to instance name
	* **value_from** - name of column whose content is used as the actual metric value

* **databases** - contains all defined databases which will be established connection, database block includes:
	* **name** - identify database block, needs to be unique
	* **driver** - database's driver ("mysql" | "postgres"),
	* **driver_option** - block which defines dns option such like hostname, port (if not given, the defaults for the driver will be set), username, password and name of database)
	* **selectdb** - name of database to which the plugin will switch after the connection is established (optional)
	* **dbqueries** - block of queries associates with this database connection

### Collected Metrics

Metric's namespace is `/intel/dbi/<metric name>/`.


Depending on the configuration, the returned values are then converted into metrics. In examples there are ready configuration setfiles with prepared queries about:
																												
[**a) Openstack Cinder services**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_cinder_services.json)
</br>[**b) Openstack Neutron agents**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_neutron_agents.json)
</br>[**c) Openstack Nova services**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_nova_services.json)
</br>[**d) Openstack Nova cluster status**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_nova_cluster_status.json)			

Plus all of above in one config file [dbi_openstack.json](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_openstack.json)


Task manifest contains names of metrics which will be collected. Also, there is posibilities of declaring `/intel/dbi/*` what means that all available metrics will be collected.

By default metrics are gathered once per second.

### Examples

Example of running snap dbi collector of and writing results to file.

Create configuration file for dbi plugin (example from [examples/configs/setfiles/dbi_openstack.json](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_openstack.json))
```json
{
    "databases": [
        {
            "name": "cinder",
            "driver": "mysql",
            "driver_option": {
                "host": "123.456.78.9",
                "port": "3306",
                "username": "cinder",
                "password": "passwd",
                "dbname": "cinder"
            },
            "dbqueries": [
                {
                    "query": "cinder_services_up"
                },
                {
                    "query": "cinder_services_down"
                },
                {
                    "query": "cinder_services_disabled"
                }
            ]
        },
        {
            "name": "neutron",
            "driver": "mysql",
            "driver_option": {
                "host": "123.456.78.9",
                "port": "3306",
                "username": "neutron",
                "password": "passwd",
                "dbname": "neutron"
            },
            "dbqueries": [
                {
                    "query": "neutron_agents_down"
                },
                {
                    "query": "neutron_agents_up"
                },
                {
                    "query": "neutron_agents_disabled"
                }
            ]
        },
        {
            "name": "nova",
            "driver": "mysql",
            "driver_option": {
                "host": "123.456.78.9",
                "port": "3306",
                "username": "nova",
                "password": "passwd",
                "dbname": "nova"
            },
            "dbqueries": [
                {
                    "query": "nova_wsrep_ready"
                },
                {
                    "query": "nova_wsrep_cluster_status"
                },
                {
                    "query": "nova_wsrep_cluster"
                },
                {
                    "query": "nova_services_down"
                },
                {
                    "query": "nova_services_up"
                },
                {
                    "query": "nova_services_disabled"
                }
            ]
        }
    ],
    "queries": [
        {
            "name": "cinder_services_down",
            "statement": "select concat_ws('/', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'down') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s1.disabled=0 and s1.deleted=0 and timestampdiff(SECOND,s1.updated_at,utc_timestamp())>120 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "cinder_services_up",
            "statement": "select concat_ws('/', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'up') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s1.disabled=0 and s1.deleted=0 and timestampdiff(SECOND,s1.updated_at,utc_timestamp())<=120 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "cinder_services_disabled",
            "statement": "select concat_ws('/', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'disabled') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s2.disabled = 1 and s1.deleted=0 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "neutron_agents_down",
            "statement": "select concat_ws('/', 'agents', replace(replace(a1.binary, '-agent', ''), 'neutron-', ''), 'down') as metric, count(a2.id) as value from agents a1 left outer join agents a2 on a1.id = a2.id and a1.admin_state_up=1 and timestampdiff(SECOND,a1.heartbeat_timestamp,utc_timestamp())>60 group by a1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "neutron_agents_up",
            "statement": "select concat_ws('/', 'agents', replace(replace(a1.binary, '-agent', ''), 'neutron-', ''), 'up') as metric, count(a2.id) as value from agents a1 left outer join agents a2 on a1.id = a2.id and a1.admin_state_up=1 and timestampdiff(SECOND,a1.heartbeat_timestamp,utc_timestamp())<60  group by a1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "neutron_agents_disabled",
            "statement": "select concat_ws('/', 'agents', replace(replace(a1.binary, '-agent', ''), 'neutron-', ''), 'disabled') as metric, count(a2.id) as value from agents a1 left outer join agents a2 on a1.id = a2.id and a1.admin_state_up=0 group by a1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "nova_wsrep_ready",
            "statement": "select replace(lower(VARIABLE_NAME), 'wsrep_', 'cluster/') as metric,  cast(replace(replace(VARIABLE_VALUE, 'ON', '1'), 'OFF', '0') as unsigned int) as value from information_schema.GLOBAL_STATUS where VARIABLE_NAME IN ('wsrep_ready', 'wsrep_connected')",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "nova_wsrep_cluster_status",
            "statement": "select replace(replace(lower(VARIABLE_NAME), 'wsrep_', ''), '_', '/') as metric, cast(replace(replace(replace(VARIABLE_VALUE, 'Primary', '1'), 'Non-Primary', '2'), 'Disconnected', '3') as unsigned int) as value from information_schema.GLOBAL_STATUS where VARIABLE_NAME = 'wsrep_cluster_status'",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "nova_wsrep_cluster",
            "statement": "select replace(replace(lower(VARIABLE_NAME), 'wsrep_', 'cluster/'), 'cluster_size', 'size') as metric, cast(VARIABLE_VALUE as unsigned int) as value from information_schema.GLOBAL_STATUS where VARIABLE_NAME IN ('wsrep_cluster_size', 'wsrep_replicated', 'wsrep_replicated_bytes', 'wsrep_received_bytes', 'wsrep_received', 'wsrep_local_commits', 'wsrep_local_cert_failures', 'wsrep_local_send_queue', 'Slow_queries')",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "nova_services_down",
            "statement": "select concat_ws('/', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'down') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s1.disabled=0 and s1.deleted=0 and timestampdiff(SECOND,s1.updated_at,utc_timestamp())>120 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "nova_services_up",
            "statement": "select concat_ws('/', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'up') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s1.disabled=0 and s1.deleted=0 and timestampdiff(SECOND,s1.updated_at,utc_timestamp())<=120 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "nova_services_disabled",
            "statement": "select concat_ws('/', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'disabled') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s2.disabled = 1 and s1.deleted=0 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        }
    ]
}
    

```

Set path to configuration file as the field `setfile` in Global Config (example from examples/configs/snap-config-sample.json):
```json

    "control": {
        "cache_ttl": "5s"
    },
    "scheduler": {
        "default_deadline": "5s",
        "worker_pool_size": 5
    },
    "plugins": {
        "collector": {
            "dbi": {
                "all": {
                  "setfile": "$SNAP_DBI_PLUGIN_DIR/examples/configs/setfiles/dbi_openstack.json"
                }
            }
        },
        "publisher": {
            "influxdb": {
                "all": {
                    "server": "xyz.local",
                    "password": "$password"
                }
            }
        },
        "processor": {}
    }
}

```

Run the snap daemon:
```
$ snapd -l 1 -t 0 --config $SNAP_DBI_PLUGIN_DIR/examples/configs/snap-config-sample.json
```

Load dbi plugin for collecting:
```
$ snapctl plugin load $SNAP_DBI_PLUGIN_DIR/build/rootfs/snap-plugin-collector-dbi
Plugin loaded
Name: dbi
Version: 1
Type: collector
Signed: false
Loaded Time: Tue, 16 Feb 2016 12:40:07 UTC
```

See available metrics:
```
$ snapctl metric list

NAMESPACE                                        VERSIONS
/intel/dbi/*                                     1
/intel/dbi/cinder/services/backup/disabled       1
/intel/dbi/cinder/services/backup/down           1
/intel/dbi/cinder/services/backup/up             1
/intel/dbi/cinder/services/scheduler/disabled    1
/intel/dbi/cinder/services/scheduler/down        1
/intel/dbi/cinder/services/scheduler/up          1
/intel/dbi/cinder/services/volume/disabled       1
/intel/dbi/cinder/services/volume/down           1
/intel/dbi/cinder/services/volume/up             1
/intel/dbi/neutron/agents/dhcp/disabled          1
/intel/dbi/neutron/agents/dhcp/down              1
/intel/dbi/neutron/agents/dhcp/up                1
/intel/dbi/neutron/agents/l3/disabled            1
/intel/dbi/neutron/agents/l3/down                1
/intel/dbi/neutron/agents/l3/up                  1
/intel/dbi/neutron/agents/metadata/disabled      1
/intel/dbi/neutron/agents/metadata/down          1
/intel/dbi/neutron/agents/metadata/up            1
/intel/dbi/neutron/agents/openvswitch/disabled   1
/intel/dbi/neutron/agents/openvswitch/down       1
/intel/dbi/neutron/agents/openvswitch/up         1
/intel/dbi/nova/cluster/connected                1
/intel/dbi/nova/cluster/local_cert_failures      1
/intel/dbi/nova/cluster/local_commits            1
/intel/dbi/nova/cluster/local_send_queue         1
/intel/dbi/nova/cluster/ready                    1
/intel/dbi/nova/cluster/received                 1
/intel/dbi/nova/cluster/received_bytes           1
/intel/dbi/nova/cluster/replicated               1
/intel/dbi/nova/cluster/replicated_bytes         1
/intel/dbi/nova/cluster/size                     1
/intel/dbi/nova/cluster/status                   1
/intel/dbi/nova/services/cert/disabled           1
/intel/dbi/nova/services/cert/down               1
/intel/dbi/nova/services/cert/up                 1
/intel/dbi/nova/services/compute/disabled        1
/intel/dbi/nova/services/compute/down            1
/intel/dbi/nova/services/compute/up              1
/intel/dbi/nova/services/conductor/disabled      1
/intel/dbi/nova/services/conductor/down          1
/intel/dbi/nova/services/conductor/up            1
/intel/dbi/nova/services/consoleauth/disabled    1
/intel/dbi/nova/services/consoleauth/down        1
/intel/dbi/nova/services/consoleauth/up          1
/intel/dbi/nova/services/scheduler/disabled      1
/intel/dbi/nova/services/scheduler/down          1
/intel/dbi/nova/services/scheduler/up            1
/intel/dbi/nova/slow_queries                     1

```

Load file plugin for publishing:
```
$ snapctl plugin load $SNAP_DIR/build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Tue, 16 Feb 2016 12:42:10 UTC
```

Create a task JSON file using for creating a task collecting all available dbi metrics (exemplary file in examples/tasks/dbi-file.json):  
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/dbi/*": {}
            },
            "config": {},
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_dbi"
                    }
                }
            ]
        }
    }
}
    
```

Create a task:
```
$ snapctl task create -t $SNAP_DBI_PLUGIN_DIR/examples/tasks/dbi-file.json
Using task manifest to create task
Task created
ID: da9188b4-d592-4b45-b108-de06a8fdee1a
Name: Task-da9188b4-d592-4b45-b108-de06a8fdee1a
State: Running
```
See sample output from `snapctl task watch <task_id>`

```
$ snapctl task watch da9188b4-d592-4b45-b108-de06a8fdee1a

Watching Task (da9188b4-d592-4b45-b108-de06a8fdee1a):
NAMESPACE                                        DATA                    TIMESTAMP                                       SOURCE
/intel/dbi/cinder/services/backup/disabled       0                       2016-02-16 12:43:12.135743129 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/backup/down           0                       2016-02-16 12:43:12.135787905 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/backup/up             1                       2016-02-16 12:43:12.135734564 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/scheduler/disabled    0                       2016-02-16 12:43:12.135790187 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/scheduler/down        0                       2016-02-16 12:43:12.135674327 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/scheduler/up          1                       2016-02-16 12:43:12.135736637 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/volume/disabled       0                       2016-02-16 12:43:12.135716351 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/volume/down           0                       2016-02-16 12:43:12.135738704 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/volume/up             1                       2016-02-16 12:43:12.135785904 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/dhcp/disabled          0                       2016-02-16 12:43:12.13572654 +0000 UTC          node-22.domain.tld
/intel/dbi/neutron/agents/dhcp/down              0                       2016-02-16 12:43:12.135758258 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/dhcp/up                1                       2016-02-16 12:43:12.13577611 +0000 UTC          node-22.domain.tld
/intel/dbi/neutron/agents/l3/disabled            0                       2016-02-16 12:43:12.135704307 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/l3/down                0                       2016-02-16 12:43:12.135728479 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/l3/up                  1                       2016-02-16 12:43:12.135769895 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/metadata/disabled      0                       2016-02-16 12:43:12.135708729 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/metadata/down          0                       2016-02-16 12:43:12.135760312 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/metadata/up            1                       2016-02-16 12:43:12.135730503 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/openvswitch/disabled   0                       2016-02-16 12:43:12.135685201 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/openvswitch/down       0                       2016-02-16 12:43:12.135702254 +0000 UTC         node-22.domain.tld
/intel/dbi/neutron/agents/openvswitch/up         3                       2016-02-16 12:43:12.135724492 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/connected                1                       2016-02-16 12:43:12.13571047 +0000 UTC          node-22.domain.tld
/intel/dbi/nova/cluster/local_cert_failures      0                       2016-02-16 12:43:12.135706467 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/local_commits            1.347723e+06            2016-02-16 12:43:12.135682691 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/local_send_queue         0                       2016-02-16 12:43:12.135751054 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/ready                    1                       2016-02-16 12:43:12.13567725 +0000 UTC          node-22.domain.tld
/intel/dbi/nova/cluster/received                 10608                   2016-02-16 12:43:12.135744873 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/received_bytes           85002                   2016-02-16 12:43:12.135778085 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/replicated               1.348857e+06            2016-02-16 12:43:12.135689414 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/replicated_bytes         1.158715258e+09         2016-02-16 12:43:12.135749002 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/size                     1                       2016-02-16 12:43:12.135752588 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/cluster/status                   1                       2016-02-16 12:43:12.135712136 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/cert/disabled           0                       2016-02-16 12:43:12.135783981 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/cert/down               0                       2016-02-16 12:43:12.135754395 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/cert/up                 1                       2016-02-16 12:43:12.135696593 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/compute/disabled        0                       2016-02-16 12:43:12.135732626 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/compute/down            0                       2016-02-16 12:43:12.135746829 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/compute/up              2                       2016-02-16 12:43:12.135781998 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/conductor/disabled      0                       2016-02-16 12:43:12.13574092 +0000 UTC          node-22.domain.tld
/intel/dbi/nova/services/conductor/down          0                       2016-02-16 12:43:12.135694604 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/conductor/up            1                       2016-02-16 12:43:12.135756342 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/consoleauth/disabled    0                       2016-02-16 12:43:12.13579255 +0000 UTC          node-22.domain.tld
/intel/dbi/nova/services/consoleauth/down        0                       2016-02-16 12:43:12.135680226 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/consoleauth/up          1                       2016-02-16 12:43:12.135772075 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/scheduler/disabled      0                       2016-02-16 12:43:12.135774278 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/scheduler/down          0                       2016-02-16 12:43:12.135780123 +0000 UTC         node-22.domain.tld
/intel/dbi/nova/services/scheduler/up            1                       2016-02-16 12:43:12.13571416 +0000 UTC          node-22.domain.tld
/intel/dbi/nova/slow_queries                     0                       2016-02-16 12:43:12.135687558 +0000 UTC         node-22.domain.tld

```
(Keys `ctrl+c` terminate task watcher)


These data are published to file and stored there (in this example in /tmp/published_dbi).

Stop task:
```
$ snapctl task stop da9188b4-d592-4b45-b108-de06a8fdee1a
Task stopped:
ID: da9188b4-d592-4b45-b108-de06a8fdee1a
```

### Roadmap

There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-dbi/issues).

## Community Support
This repository is one of **many** plugins in the **Snap Framework**: a powerful telemetry agent framework. To reach out on other use cases, visit:

* [Snap Gitter channel] (https://gitter.im/intelsdi-x/snap)

The full project is at http://github.com:intelsdi-x/snap.

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: 	[Izabella Raulin](https://github.com/IzabellaRaulin)

**Thank you!** Your contribution is incredibly important to us.
