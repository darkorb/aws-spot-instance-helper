aws-spot-instance-helper
========

A microservice that does micro things.

## Building

`make`


## Running

`./bin/aws-spot-instance-helper`

## Usage and env variables

```
--debug, -d                             
        Debug logging [$DEBUG]

--poll-interval value, -i value         
        Polling interval for checks (default: 5s) [$POLL_INTERVAL]

--cattleURL value, -u value             
        Cattle URL [$CATTLE_URL]

--cattleAccessKey value, --ck value     
        Cattle Access Key [$CATTLE_ACCESS_KEY]

--cattleSecretKey value, --cs value     
        Cattle Secret Key [$CATTLE_SECRET_KEY]

--slackWebhookUrl value, -s value       
        Slack Webhook URL [$SLACK_WEBHOOK]

--slackMessageSuffix value, --ss value  
        Appears at the end of slack message - eg: @bob.white or mytracker63763 etc [$SLACK_MESSAGE_SUFFIX]

--slackInitAnnouncement, --si boolean        
        Initial announcement will be send to slack on a startup [$SLACK_INIT_ANNOUNCEMENT]

--help, -h                              
        show help

--version, -v                           
        print the version
```

## License
Copyright (c) 2014-2016 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
