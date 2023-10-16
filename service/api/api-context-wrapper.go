package api

// Il metodo wrap è utilizzato per avvolgere gli handler HTTP, fornendo un contesto di richiesta che include un UUID univoco
// e un logger configurato. Questo permette agli handler HTTP di avere accesso a informazioni contestuali sulla richiesta e
// un logger che automaticamente include informazioni utili nei messaggi di log.
 
import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
// il wrap serve a creare un nuovo handler HTTP che "avvolge" fn, fornendo il contesto della richiesta come parametro aggiuntivo.
// permettendo agli handler HTTP di avere accesso a informazioni contestuali sulla richiesta.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {

	// La funzione restituita da wrap genera un nuovo UUID per ogni richiesta,
	// crea un nuovo reqcontext.RequestContext, e chiama fn con il contesto della richiesta come parametro aggiuntivo.
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()

		//Se la generazione dell'UUID fallisce, logga un errore e restituisce un errore HTTP 500 al client.
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// ctx.Logger è configurato per includere l'UUID della richiesta e l'indirizzo IP remoto in tutti i messaggi di log.
		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		// fn (la funzione handler originale) viene chiamata con il contesto della richiesta come ultimo parametro.
		fn(w, r, ps, ctx)
	}
}
