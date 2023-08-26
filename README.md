# YAML -> JSON

## Goal

Convert YAML files to JSON files in a streaming fashion for performance.

This is a toy project for me to learn about Go channels.

YAML spec: https://yaml.org/spec/1.2.2/

Thanks to these online tools:

- https://jsonformatter.org/yaml-to-json
- https://www.yamllint.com/

## Limitations

As this is just a prototype, there are many limitations.

- only a small part of the YAML spec is implemented
- error handling has not been given much thought

## Example

Input:

```yaml
data:
  - name: John
    age: 30
  - name: Jane
    age: 25
```

Output (pretty-printed):

```json
{
  "data": [
    {
      "name": "John",
      "age": 30
    },
    {
      "name": "Jane",
      "age": 25
    }
  ]
}
```
