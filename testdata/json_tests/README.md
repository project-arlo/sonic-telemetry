JSON Schema Based Tests
=======================

In order to validate the JSON output of gNMI and REST interfaces in a
flexible way that allows the same test cases to work on multiple
platforms, JSON Schema is proposed. The features of JSON Schema that are
useful for this purpose are:

-   Structure Validation: Validate the response is in the right format
    (objects, arrays, field types etc)
-   Data Validation: Validate the returned data is what we expect,
    using:

    -   Exact String Match
    -   Numerical range match
    -   Regular Expression String Matching
    -   Enum matching
    -   List content matching, even if list contains complex objects by
        using a sub-schema.

-   Single schema file can contain all the data needed to run a test
    since we can add extra fields that we use in our code, but are
    ignored by JSON Schema validation library
-   Existing JSON Schema validation libraries for go, python C/C++ and
    more.
-   When fields or structure don’t match, you get an exact error
    matching telling you what the issue is in the test results, instead
    of an opaque pass/fail error.

Example of GET Test
-------------------
```
{
  "operation": "get",
  "returnCode": 0,
  "xpath": "openconfig-interfaces:interfaces/interface[name=Ethernet0]",
  "target": "OC_YANG",

  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "title": "Check Interface mtu",
  "required": [
    "openconfig-interfaces:interface"
  ],
  "properties": {
    "openconfig-interfaces:interface": {
      "type": "array",
      "title": "The Openconfig-interfaces:interface Schema",
      "items": {
        "type": "object",
        "title": "The Items Schema",
        "required": [
          "config"
        ],
        "properties": {
          "config": {
            "type": "object",
            "title": "The Config Schema",
            "required": [
              "mtu"
            ],
            "properties": {
              "mtu": {
                "type": "integer",
                "title": "The Mtu Schema",
                "enum": [
                  9003
                ]
              }
            }
          }
        }
      }
    }
  }
}
```
### Test Explanation

The first five fields are custom fields that I inserted to instruct the
go test function what operation is being performed, the expected gNMI
return code, the gNMI path to get and the DB/target it is directed at.
JSON Schema validator will ignore these five fields.

The remaining fields are JSON Schema specific fields that you can read
about here: <https://json-schema.org/specification.html>. However the
title field in the top level object is re-used as the title of the go
test as well. The title and the \$id field are optional.

Example of SET Test
-------------------
```
{
  "operation": "replace",
  "returnCode": 0,
  "xpath": "openconfig-interfaces:interfaces/interface[name=Ethernet0]",
  "target": "OC_YANG",
  "title": "Set Interface mtu",
  "attributeData": "{\"config\": {\"mtu\":9003}}"

}
```

### Test Explanation

This test contains the same five fields but with the addition of
attributeData, which contains the JSON payload to send in the set
request.

There is no JSON Schema fields here since the SET request does not
return any JSON result, only a response code.
