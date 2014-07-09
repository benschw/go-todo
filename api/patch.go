package api

/**
 * on patching: http://williamdurand.fr/2014/02/14/please-do-not-patch-like-an-idiot/
 *
 * patch specification https://tools.ietf.org/html/rfc5789
 * json definition http://tools.ietf.org/html/rfc6902
 */

type Patch struct {
	Op    string `json:"op" binding:"required"`
	From  string `json:"from"`
	Path  string `json:"path"`
	Value string `json:"value"`
}
