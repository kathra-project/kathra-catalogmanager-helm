/*
 * Kathra Catalog Manager
 *
 * KATHRA Catalog Management API permetting :   * Generate source's packages from templates   * Insert catalog entry from template      * Insert catalog entry from file     * Insert catalog entry from source repository    * Read catalog entries from catalog    * Read catalog details from catalog
 *
 * API version: 1.1.0-RC-SNAPSHOT
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package kathracatalogmanagerhelmservices

type RuntimeEnvironment struct {

	// RuntimeEnvironment identifier
	Id string `json:"id,omitempty"`

	// RuntimeEnvironment name
	Name string `json:"name,omitempty"`

	// RuntimeEnvironment providerId
	ProviderId string `json:"providerId,omitempty"`

	// Runtime instance
	CatalogEntryInstances []CatalogEntryInstance `json:"catalogEntryInstances,omitempty"`
}
