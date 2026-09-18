import "@typespec/http";
import "../common/parameters.tsp";

namespace Shared;

/**
 * Resource filters.
 */
@friendlyName("ResourceFilters")
model ResourceFilters {
  #suppress "@meterforge/api-spec-aip/doc-decorator" "shared model"
  name?: Common.StringFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "shared model"
  labels?: Common.LabelsFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "shared model"
  public_labels?: Common.LabelsFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "shared model"
  created_at?: Common.DateTimeFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "shared model"
  updated_at?: Common.DateTimeFieldFilter;
  #suppress "@meterforge/api-spec-aip/doc-decorator" "shared model"
  deleted_at?: Common.DateTimeFieldFilter;
}
