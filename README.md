# Snap collector plugin - dbi

This plugin connects to various databases, executes SQL statements and reads back the results. 

The plugin is a generic plugin. You can configure how each column is to be interpreted and the plugin will generate one or more data sets from each row returned according. These rules are defined in separated json file (read more about this in [Configuration and Usage](#configuration-and-usage) section).

Depending on the configuration, the returned values are converted into metrics.

The plugin is used in the [Snap framework] (http://github.com/intelsdi-x/snap).				

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
#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-dbi/releases) page. Download the plugin from the latest release and load it into `snapteld` (`/opt/snap/plugins` is the default location for Snap packages).

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-dbi  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-dbi.git
```

Build the Snap dbi plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage

* Set up the [Snap framework](https://github.com/intelsdi-x/snap#getting-started)

* Create configuration file (called as a setfile) in which will be defined databases, queries and rules how interpret the results, see exemplary in [examples/configs/setfiles](examples/configs/setfiles)

* Set up field `setfile` in Global Config as a path to dbi plugin configuration file, see exemplary Snap Global Config: in [examples/configs/snap-config-sample.json] (examples/configs/snap-config-sample.json)
 
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

Metric's namespace is `/intel/dbi/<metric_name>/`.


Depending on the configuration, the returned values are then converted into metrics. In examples there are ready configuration setfiles with prepared queries about:
																												
[**a) Openstack Cinder services**](examples/configs/setfiles/dbi_cinder_services.json)
</br>[**b) Openstack Neutron agents**](examples/configs/setfiles/dbi_neutron_agents.json)
</br>[**c) Openstack Nova services**](examples/configs/setfiles/dbi_nova_services.json)
</br>[**d) Openstack Nova cluster status**](examples/configs/setfiles/dbi_nova_cluster_status.json)

Plus all of above in one config file [dbi_openstack.json](examples/configs/setfiles/dbi_openstack.json)


Task manifest contains names of metrics which will be collected

By default metrics are gathered once per second.

### Example
Example of running Snap dbi collector retrieving the metrics about cinder services and writing the results to file.

Create configuration file (`setfile`) for dbi plugin (see [examples/configs/setfiles/dbi_cinder_services.json](examples/configs/setfiles/dbi_cinder_services.json)).
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
        }
    ]
}
    

```


Create global configuration file for Snap and set path to `setfile` (see [examples/configs/snap-config-sample.json](examples/configs/snap-config-sample.json)):
```json
{
    "control": {
        "plugins": {
            "collector": {
                "dbi": {
                    "all": {
                        "setfile": "/path/to/dbi_cinder_services.json"
                    }
                }
            },
            "publisher": {},
            "processor": {}
        }
    }
}

```
Run the Snap daemon with created global config:
```
$ snapteld -l 1 -t 0 --config snap-config-sample.json
```

Download and load dbi plugin for collecting and file publisher:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-dbi/latest/linux/x86_64/snap-plugin-collector-dbi
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snaptel plugin load snap-plugin-collector-dbi
Plugin loaded
Name: dbi
Version: 4
Type: collector
Signed: false
Loaded Time: Tue, 05 Apr 2016 07:49:53 UTC
$ snaptel plugin load snap-plugin-publisher-file
Plugin loaded
Name: file
Version: 4
Type: publisher
Signed: false
Loaded Time: Tue, 05 Apr 2016 07:49:54 UTC
```

See available metrics:
```
$ snaptel metric list

NAMESPACE                                        VERSIONS
/intel/dbi/cinder/services/backup/disabled       4
/intel/dbi/cinder/services/backup/down           4
/intel/dbi/cinder/services/backup/up             4
/intel/dbi/cinder/services/scheduler/disabled    4
/intel/dbi/cinder/services/scheduler/down        4
/intel/dbi/cinder/services/scheduler/up          4
/intel/dbi/cinder/services/volume/disabled       4
/intel/dbi/cinder/services/volume/down           4
/intel/dbi/cinder/services/volume/up             4
```

Create and load example task file, which describes collection of cinder services dbi metrics (see [examples/tasks/dbi_cinder_services-file.json](examples/tasks/dbi_cinder_services-file.json)):
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
                "/intel/dbi/cinder/services/backup/disabled": {},
                "/intel/dbi/cinder/services/backup/down": {},
                "/intel/dbi/cinder/services/backup/up": {},
                "/intel/dbi/cinder/services/scheduler/disabled": {},
                "/intel/dbi/cinder/services/scheduler/down": {},
                "/intel/dbi/cinder/services/scheduler/up": {},
                "/intel/dbi/cinder/services/volume/disabled": {},
                "/intel/dbi/cinder/services/volume/down": {},
                "/intel/dbi/cinder/services/volume/up": {}
            },
            "config": {},
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_dbi_cinder_services.txt"
                    }
                }
            ]
        }
    }
}
    
```

```
$ snaptel task create -t dbi_cinder_services-file.json
Using task manifest to create task
Task created
ID: da9188b4-d592-4b45-b108-de06a8fdee1a
Name: Task-da9188b4-d592-4b45-b108-de06a8fdee1a
State: Running
```
See sample output from `snaptel task watch <task_id>`

```
$ snaptel task watch da9188b4-d592-4b45-b108-de06a8fdee1a

Watching Task (da9188b4-d592-4b45-b108-de06a8fdee1a):
NAMESPACE                                        DATA                    TIMESTAMP                                       SOURCE
/intel/dbi/cinder/services/backup/disabled       0                       2016-04-05 08:43:12.135743129 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/backup/down           0                       2016-04-05 08:43:12.135787905 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/backup/up             1                       2016-04-05 08:43:12.135734564 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/scheduler/disabled    0                       2016-04-05 08:43:12.135790187 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/scheduler/down        0                       2016-04-05 08:43:12.135674327 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/scheduler/up          1                       2016-04-05 08:43:12.135736637 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/volume/disabled       0                       2016-04-05 08:43:12.135716351 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/volume/down           0                       2016-04-05 08:43:12.135738704 +0000 UTC         node-22.domain.tld
/intel/dbi/cinder/services/volume/up             1                       2016-04-05 08:43:12.135785904 +0000 UTC         node-22.domain.tld
```
(`ctrl+c` terminates task watcher)


Metrics are published to file (in this example in /tmp/published_dbi_cinder_services.txt).

Stopping task:
```
$ snaptel task stop da9188b4-d592-4b45-b108-de06a8fdee1a
Task stopped:
ID: da9188b4-d592-4b45-b108-de06a8fdee1a
```

### Roadmap

There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-dbi/issues).

## Community Support
This repository is one of **many** plugins in the **Snap Framework**: a powerful telemetry agent framework. The full project is at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: 	[Izabella Raulin](https://github.com/IzabellaRaulin)

**Thank you!** Your contribution is incredibly important to us.
