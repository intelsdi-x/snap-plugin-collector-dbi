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
- Access to database (currently following sql drivers are supported: **MySQL**, **PostgreSQL**)

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

* Create configuration file (called as a setfile) in which will be defined databases, queries and rules how interpret the results (see examplary in [examples/configs/setfiles](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/master/examples/configs/setfiles))

* Set up field `setfile` in Global Config as a path to dbi plugin configuration file. Sample Global Config is available in [examples/configs/snap-config-sample.json] (https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/master/examples/configs/snap-config-sample.json)
 
Notice that this plugin is a generic plugin, i.e. it cannot work without configuration, because there is no reasonable default behavior.

## Documentation

### Setfile fields
* **queries** - contains all defined queries put in query block which includes:
	*  **name** - identify query block,
	*  **statement** - SQL statement to be executed
	*  **results** - block which defines results of statement
* **results** - contains how the returned data should be interpreted, including:
	 * **name** - name of result, acceptable empty if only one result is defined; in other case must be given in order to distinguish results
	* **instance_from** - name of column whose values will be used to specify an instance
	* **instance_prefix** - prepended prefix to instance name
	* **value_from** - name of column whose content is used as the actual metric value

* **databases** - contains all defined databases which will be established connection, database block includes:
	* **name** - identify database block,
	* **driver** - database's driver ("mysql" | "postgres"),
	* **driver_option** - block which defines dns option such like hostname, port (if not given, the defaults for the driver will be set), username, password and name of database)
	* **selectdb** - name of database to which the plugin will switch after the connection is established (optional)
	* **dbqueries** - block of queries associates with this database connection

### Collected Metrics

Metric's namespace is `/intel/dbi/<metric name>/`. 

Depending on the configuration, the returned values are then converted into metrics. In examples there are ready configuration setfiles with prepared queries about:																														[**a) Openstack Cinder services**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_cinder_services.json)																		[**b) Openstack Neutron agents**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_neutron_agents.json)																		[**c) Openstack Nova services**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_nova_services.json)																		[**d) Openstack Nova cluster status**](https://github.com/intelsdi-x/snap-plugin-collector-dbi/blob/dbi_plugin/examples/configs/setfiles/dbi_nova_cluster_status.json)

Also, there is posibilities of declaring in task manifest metric `/intel/dbi/*` what means that all available metrics will be collected.

By default metrics are gathered once per second.

### Examples

Example of running snap dbi collector of and writing results to file.

Create configuration file for dbi plugin (example from examples/configs/setfiles/dbi_cinder_services.json)
```json
{
    "queries": [
        {
            "name": "services_down",
            "statement": "select concat_ws('.', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'down') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s1.disabled=0 and s1.deleted=0 and timestampdiff(SECOND,s1.updated_at,utc_timestamp())>120 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "services_up",
            "statement": "select concat_ws('.', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'up') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s1.disabled=0 and s1.deleted=0 and timestampdiff(SECOND,s1.updated_at,utc_timestamp())<=120 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        },
        {
            "name": "services_disabled",
            "statement": "select concat_ws('.', 'services', replace(replace(s1.binary, 'nova-', ''), 'cinder-', ''), 'disabled') as metric, count(s2.id) as value from services s1 left outer join services s2 on s1.id = s2.id and s2.disabled = 1 and s1.deleted=0 group by s1.binary",
            "results": [
                {
                    "instance_from": "metric",
                    "value_from": "value"
                }
            ]
        }
    ],
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
                    "query": "services_up"
                },
                {
                    "query": "services_down"
                },
                {
                    "query": "services_disabled"
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
                  "setfile": "$SNAP_DBI_PLUGIN_DIR/examples/configs/setfiles/dbi_cinder_services.json"
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
Loaded Time: Tue, 12 Jan 2016 05:25:35 EST
```

See available metrics:
```
$ snapctl metric list

NAMESPACE                                                             	VERSIONS
/intel/dbi/*                                                             1
/intel/dbi/cinder/services_disabled/services.backup.disabled             1
/intel/dbi/cinder/services_disabled/services.scheduler.disabled          1
/intel/dbi/cinder/services_disabled/services.volume.disabled             1
/intel/dbi/cinder/services_down/services.backup.down                     1
/intel/dbi/cinder/services_down/services.scheduler.down                  1
/intel/dbi/cinder/services_down/services.volume.down                     1
/intel/dbi/cinder/services_up/services.backup.up                         1
/intel/dbi/cinder/services_up/services.scheduler.up                      1
/intel/dbi/cinder/services_up/services.volume.up                         1

```

Load file plugin for publishing:
```
$ snapctl plugin load $SNAP_DIR/build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Tue, 12 Jan 2016 05:26:21 EST
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
NAMESPACE                                                         DATA    TIMESTAMP                                  SOURCE
/intel/dbi/cinder/services_disabled/services.backup.disabled      0       2016-02-12 10:48:15.048385186 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_disabled/services.scheduler.disabled   0       2016-02-12 10:48:15.048394479 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_disabled/services.volume.disabled      0       2016-02-12 10:48:15.048402338 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_down/services.backup.down              0       2016-02-12 10:48:15.048439565 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_down/services.scheduler.down           0       2016-02-12 10:48:15.048410013 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_down/services.volume.down              0       2016-02-12 10:48:15.048415811 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_up/services.backup.up                  1       2016-02-12 10:48:15.048422497 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_up/services.scheduler.up               1       2016-02-12 10:48:15.048428271 +0000 UTC    node-22.domain.tld
/intel/dbi/cinder/services_up/services.volume.up                  1       2016-02-12 10:48:15.048433711 +0000 UTC    node-22.domain.tld
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
