package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func expandAzureRmCosmosDBIndexingPolicyIncludedPaths(input []interface{}) *[]documentdb.IncludedPath {
	if len(input) == 0 {
		return nil
	}

	var includedPaths []documentdb.IncludedPath

	for _, v := range input {
		includedPath := v.(map[string]interface{})
		path := documentdb.IncludedPath{
			Path: utils.String(includedPath["path"].(string)),
		}

		if indexes, ok := includedPath["index"].([]interface{}); ok {
			path.Indexes = expandCosmosDBIndexingPolicyIncludedPathIndexes(indexes)
		}

		includedPaths = append(includedPaths, path)
	}

	return &includedPaths
}

func expandCosmosDBIndexingPolicyIncludedPathIndexes(input []interface{}) *[]documentdb.Indexes {
	if len(input) == 0 {
		return nil
	}
	var indexes []documentdb.Indexes

	for _, i := range input {
		index := i.(map[string]interface{})
		indexes = append(indexes, documentdb.Indexes{
			DataType:  documentdb.DataType(index["data_type"].(string)),
			Precision: utils.Int32(int32(index["precision"].(int))),
			Kind:      documentdb.IndexKind(index["kind"].(string)),
		})
	}

	return &indexes
}

func expandAzureRmCosmosDBIndexingPolicyExcludedPaths(input []interface{}) *[]documentdb.ExcludedPath {
	if len(input) == 0 {
		return nil
	}

	var paths []documentdb.ExcludedPath

	for _, v := range input {
		block := v.(map[string]interface{})
		paths = append(paths, documentdb.ExcludedPath{
			Path: utils.String(block["path"].(string)),
		})
	}

	return &paths
}

func ExpandAzureRmCosmosDbIndexingPolicy(d *schema.ResourceData) *documentdb.IndexingPolicy {
	i := d.Get("indexing_policy").([]interface{})
	policy := &documentdb.IndexingPolicy{
		Automatic: utils.Bool(true),
	}

	if len(i) == 0 || i[0] == nil {
		policy.IndexingMode = documentdb.Consistent
		policy.IncludedPaths = &[]documentdb.IncludedPath{
			{Path: utils.String("/*")},
		}
		policy.ExcludedPaths = &[]documentdb.ExcludedPath{}

		return policy
	}

	input := i[0].(map[string]interface{})
	policy.IndexingMode = documentdb.IndexingMode(input["indexing_mode"].(string))
	if v, ok := input["included_path"].([]interface{}); ok {
		policy.IncludedPaths = expandAzureRmCosmosDBIndexingPolicyIncludedPaths(v)
	}
	if v, ok := input["excluded_path"].([]interface{}); ok {
		policy.ExcludedPaths = expandAzureRmCosmosDBIndexingPolicyExcludedPaths(v)
	}

	return policy
}

func flattenCosmosDBIndexingPolicyExcludedPaths(input *[]documentdb.ExcludedPath) []interface{} {
	if input == nil {
		return nil
	}

	excludedPaths := make([]interface{}, 0)

	for _, v := range *input {
		// _etag is automatically added by the server and should be excluded on flattening
		// as the user isn't setting it and it will show changes in state.
		if *v.Path == "/\"_etag\"/?" {
			continue
		}

		block := make(map[string]interface{})
		block["path"] = v.Path
		excludedPaths = append(excludedPaths, block)
	}

	return excludedPaths
}

func flattenCosmosDBIndexingPolicyIncludedPaths(input *[]documentdb.IncludedPath) []interface{} {
	if input == nil {
		return nil
	}

	includedPaths := make([]interface{}, 0)

	for _, v := range *input {
		block := make(map[string]interface{})
		block["path"] = v.Path
		block["index"] = flattenCosmosDBIndexingPolicyIncludedPathIndexes(v.Indexes)
		includedPaths = append(includedPaths, block)
	}

	return includedPaths
}

func flattenCosmosDBIndexingPolicyIncludedPathIndexes(input *[]documentdb.Indexes) []interface{} {
	if input == nil {
		return nil
	}

	indexesBlocks := make([]interface{}, 0)

	for _, v := range *input {
		block := make(map[string]interface{})

		block["data_type"] = v.DataType
		block["precision"] = v.Precision
		block["kind"] = v.Kind

		indexesBlocks = append(indexesBlocks, block)
	}

	return indexesBlocks
}

func FlattenAzureRmCosmosDbIndexingPolicy(indexingPolicy *documentdb.IndexingPolicy) []interface{} {
	results := make([]interface{}, 0)
	if indexingPolicy == nil {
		return results
	}

	result := make(map[string]interface{}, 0)
	result["indexing_mode"] = string(indexingPolicy.IndexingMode)
	result["included_path"] = flattenCosmosDBIndexingPolicyIncludedPaths(indexingPolicy.IncludedPaths)
	result["excluded_path"] = flattenCosmosDBIndexingPolicyExcludedPaths(indexingPolicy.ExcludedPaths)

	results = append(results, result)
	return results
}
