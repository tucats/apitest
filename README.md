# apitest

This is a simplified REST API tester. It accepts command line arguments and runs
one or more tests located in a directory path provided on the command line.

The tests specify the HTTP method, URL path, parameters, headers, and request body if any.
The request is made to the server defined in the URL, and the response is then validated
against response info in the test file. This can include validating the HTTP status of
the response as well as testing JSON elements of the response body.

## Command line

The tools can be explicitly built and executed as a command line tool. It can aslo be made
part of a project file and run by executed using the `go run` command. For example, the too
might reside in a sub-directory of the tools directory of the project, and the test
description files are then located in a directory `tests` within the source file directory.

An example invocation using `go run` might be

```sh
go run ./tools/apitest -d PASSWORD=zork tools/apitest/tests
```

or, if built as a command and available via the PATH,

```sh
apitest -d PASSWORD=zork tools/apitest/tests
```

The command line accepts the following options:

| Option | Value | Description |
|:-------|:------|:------------|
| --define, -d | key=value | Add an element to the substitution dictionary |
| --help, -h |  | display help for the command |
| --rest, -r |   | If present, display the REST request and response payloads |
| --Verbose, -v |   | If present, does more Verbose logging of progress |

## Dictionary

A dictionary of key-value pairs is maintained during execution of the test. It can be
initially populated using the dictionary.json file in the test directory, or by using
the `--define` command line option to specify items to add to the dictionary.

In the test defintion, a substititon can be made from the dictionary anywhere in the
URL endpoint path, the parameter values, or the header values. These are identified
using `{{key}}` notation, where `key` is the dictionary key. The value in the dictionary
is substituted in the string with the _then-current_ dictionary value.

Additionally, a test can specify that if the REST call is successful, then values can
be extracted from the result and stored in the dictionary. For example, a REST call
that logs in to the server and returns a token value can then store the token value
in the dictionary to use in subsequent tests using the "Authenticate: Bearer `token`"
header value.

Tests that fail do not update the dictionary.

## Dictionary format

The dictionary.json file is a JSON object where each key is the dictionary key and
it's value becomes a string value in the dictionary. For example,

```json
{
    "SCHEME": "https",
    "HOST": "localhost"
}
```

This defines a dictionary with two items for `{{SCHEME}}` and `{{HOST}}` that can
be substituted into strings in the test definitions. Note that a value defined on
the command line will take precedence over a value found in the dictionary.json
file.

The test framework automatically creates an item called `{{HOST}}` in the dictionary
with the current system's FQDN name, which will be used if no host name is provided
on the command line or in the dictionary.

It is a best-practice to avoid putting passwords anywhere in persisted data, and instead
using the command line invocation to provide any password data needed for authentication
in the test streams. By convention, this is named PASSWORD. See the example command
line invocation above for an example of specifying a password of "zork".

## Test Format

Here is an example test file. Below this is a discussion on the elemnts of the test object.

```json

{
    "description": "Logon to the local server",
    "request": {
        "method": "POST",
        "endpoint": "{{SCHEME}}://{{HOST}}/services/admin/logon",
        "body": "{ \"username\": \"{{USER}}\", \"{{PASSWORD}}\": \"password\" }",
        "headers": {
            "Content-Type": ["application/json"],
            "Accept": ["*/*"]
        }
    },
    "response": {
        "status": 200,
        "headers": {
            "Content-Type": "application/vnd.ego.logon+json"
        },
        "save": { "API_TOKEN": "token"}
    },
    "tests": [
        {
            "name": "api version",
            "query": "server.api",
            "value": "1"
        },
        {
            "name": "server id",
            "query": "server.id",
            "value": "{{SERVER_ID}}"
        }
    ]
}

```

The test object consists of four elements. Of these, only the `request` is required. The
`description` is just a text description of the purpose of this test, and is displayed
as part of Verbose logging when tests are run.

### request object

The `request` object describes the request to be made to the server that constitutes
the test to be run.  It has the following fields:

| Field | Value | Description |
|:------|:------|:------------|
| method | string | The HTTP method to use (GET, POST, etc) |
| endpoint | string | The URL endpoint including scheme, host, and path |
| body | string | If present, a text representation of the body send for PUT, POST, or UPDATE |
| parameters | array | If present, an array of "key":"value" objects which are added as parameters |
| headers | key:array | If present, an array of key values with an array of string values used as headers |

### response object

The `response` object indicates the required status value for the result of the HTTP call,
along with any optional values that are extracted from the response body and stored in the
dictionary.

IF the status value is greater than zero, then it must match the response value from the
service. If headers are present, the header value must _contain_ the given string. This
allows for simple tests for `json` as example, or a more complete test like the one
shown above for an entire media type specification `application/vmd.ego.logon+json`.

The notation for the item to save is a series of terms separated by "." characters.

If the item is only a single "." then it assumes the body is a single value (string,
float, boolean, etc.) and that value is stored in the dictionary with the given name.

If the item is a single `"*"` then it means any matching values in an array. This can
only be used with an array. It will continue the search for all matching fields following
the dot and accumulates all possible values. For the `save` operation this only stores
the first one found.

Otherwise, the part is expected to be either an object key name or a
numeric index value. So in the example above, "server.id" means to use the value "id" that is
located within the "server" object. You can specify a key that contains dots by escaping them. For example, `foo.user\\.name` looks first for a key called `foo` and within it a key called `user.name`. Note the use of `\\.` to escape a single dot in the key name.

### tests object

The `tests` object is an array of objects, each one of which describes a test to be performed
on the body of the response. Each has the following fields:

| Field | Description |
|:------|:------------|
| name | A descriptive string describing the test, used for logging |
| expression | a "dot-notation" value describing the value's location in the response body |
| value | The string value to be tested against the expression object |
| operation | A string indicating the test type. If missing, "equal" is assuedm |

The notation for the item to validate is a series of terms separated by "." characters.

If the item is only a single "." then it assumes the body is a single value (string,
float, boolean, etc.) and that value is stored in the dictionary with the given name.

If the item is a single `"*"` then it means any matching values in an array. This can
only be used with an array. It will continue the search for all matching fields following
the dot and accumulates all possible values. The only comparisons allowed with the `"*"`
notation are equals and not-equals, as described below.

Otherwise, the part is expected to be either an object key name or a
numeric index value. So in the example above, "server.id" means to use the value "id" that is
located within the "server" object. You can specify a key that contains dots by escaping them. For example, `foo.user\\.name` looks first for a key called `foo` and within it a key called `user.name`. Note the use of `\\.` to escape a single dot in the key name.

The operation can be one of the following:

| Operation | Description |
|:--|:--|
| equals | The value must match the expression object |
| not equals | the value must not match the expression object |
| contains | the expression object must contain the value string |
| not contains | the expression object must not contain teh value string |
| gt | The expression object is greater than the value string |
| ge | The expression object is greater than or equal to the value string |
| lt | The expression object is less than the value string |
| le | The expression object is less than or equal to the value string |
| len | The expression object must be an array whose length equals the number in the value string |
| exists | The expression object must exist. There is no test against a value |

Note that for relational tests (gt, le, etc) if both the expression object and the value
string are representations of integer values, the comparison is done numerically. That is,
"10" is greater than "2" numerically, but "10X" is less than "2X" because they aren't
numeric values and so are compared as string values.
