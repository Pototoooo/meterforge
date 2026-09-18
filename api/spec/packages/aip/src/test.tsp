import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/openapi3";
import "./common/parameters.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

@service(#{ title: "Filter Test API" })
@info(#{ version: "1.0.0" })
namespace Test;

/**
 * Field filters with all supported types.
 */
@friendlyName("FieldFilters")
model FieldFilters {
  #suppress "@meterforge/api-spec-aip/doc-decorator" "test model"
  boolean?: Common.BooleanFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "test model"
  numeric?: Common.NumericFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "test model"
  string?: Common.StringFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "test model"
  string_exact?: Common.StringFieldFilterExact;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "test model"
  ulid?: Common.ULIDFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "test model"
  datetime?: Common.DateTimeFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "test model"
  labels?: Common.LabelsFieldFilter;
}

@route("/field-filters")
interface FieldFiltersEndpoints {
  @get
  @operationId("get-field-filters-with-all-supported-types")
  @summary("Field filters")
  getFieldFilters(
    @query(#{ style: "deepObject", explode: true })
    filter?: FieldFilters,
  ): void;
}
