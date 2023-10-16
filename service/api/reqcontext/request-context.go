/*
Package reqcontext contains the request context. Each request will have its own instance of RequestContext filled by the
middleware code in the api-context-wrapper.go (parent package).
Quando una richiesta HTTP arriva al server, api-context-wrapper.go crea un'istanza di RequestContext,
genera un UUID univoco per la richiesta, e configura il logger per includere questo UUID nei messaggi di log.

L'istanza di RequestContext viene poi passata, rendendo l'UUID e il logger
facilmente accessibili a tutte le parti del codice che gestiscono la richiesta.
Questo permette di tracciare facilmente l'elaborazione di una singola richiesta attraverso il sistema e nei log,
il che è particolarmente utile per il debugging e la tracciabilità in ambienti di produzione.

Each value here should be assumed valid only per request only, with some exceptions like the logger.
*/
package reqcontext

//  importo una libreria utilizzata per lavorare con UUID (Universal Unique Identifiers) ed una libreria di logging.

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// RequestContext is the context of the request, for request-dependent parameters
// È una struttura che contiene informazioni contestuali per una richiesta HTTP.
type RequestContext struct {
	// ReqUUID is the request unique ID
	ReqUUID uuid.UUID

	// Logger is a custom field logger for the request
	// Un logger che può essere utilizzato per registrare messaggi relativi a questa richiesta specifica
	Logger logrus.FieldLogger
}
