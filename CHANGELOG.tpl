{{#releases}}
## [{{name}}](https://github.com/konstellation-io/kai-kli/tag/{{name}}) ({{date}})

{{#sections}}
### {{name}}

{{#commits}}
* [[{{#short5}}{{sha}}{{/short5}}](https://github.com/konstellation-io/kai-kli/commit/{{sha}})] {{message.fullMessage}} ({{authorAction.identity.name}}, {{#timestampISO8601}}{{commitAction.timeStamp.timeStamp}}{{/timestampISO8601}})

{{/commits}}
{{^commits}}
No changes.
{{/commits}}
{{/sections}}
{{^sections}}
No changes.
{{/sections}}
{{/releases}}
{{^releases}}
No releases.
{{/releases}}
