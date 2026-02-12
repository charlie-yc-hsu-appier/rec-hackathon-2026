You are an AI coding assistant.

Read the current vendor configuration from {{VENDOR_CONFIG_FILE_PATH}}
Find the vendor named "{{VENDOR_NAME}}" and update its configuration to match the following JSON:
{{VENDOR_JSON}}

Note that the required fields of the vendor configuration are "name", "with_proxy", "http_method", "request.url", "request.queries" and "tracking.url". If any of these fields are not applicable, still keep them in the config with empty value.
Note that "body", "response_body" and "headers" are not required fields in the vendor configuration, but they are important for the implementation.

If the request method is GET, you can ignore the "body". If there are new fields in queries, check if we can support the macro for the value, if not, add the support in the code base.

If the request method is POST, you need to implement or update the "body" based on the provided JSON.
Create a file to structure the POST body and name it as {{VENDOR_NAME}}.go, if it does not exist.
Assign this body strategy to vendor "{{VENDOR_NAME}}" in function BuildBody of file {{BUILDER_FILE_PATH}}.

Read the current repo and find if there exists unmarshaler that already support parsing the above response_body, if no, create it. Assign this unmarshaler to vendor "{{VENDOR_NAME}}" in function BuildUnmarshaler of file {{BUILDER_FILE_PATH}}.

Update the request header if there is any change. Create a file for header strategy if needed and assign it to the vendor in function BuildHeader of file {{BUILDER_FILE_PATH}}.

Create or update the unit tests for the above changes. If there is just config update, no need to create or update unit tests.
          
Finally, output the updated vendor configuration file content only. Make sure the output is a valid YAML content.
