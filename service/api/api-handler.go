package api

/*
definisce un insieme di endpoint per un'API web utilizzando un router HTTP.
Ogni endpoint è associato a un metodo HTTP specifico e a un "handler" che gestisce le richieste a quell'endpoint.
Handler serve a configurare un router HTTP che sa come gestire diverse richieste HTTP a vari endpoint dell'API web,
indirizzando ciascuna richiesta all'handler appropriato e assicurando che la logica dell'applicazione sia eseguita correttamente.
    Ricezione della Richiesta:
        Quando una richiesta HTTP arriva al server, il Handler determina quale funzione handler dovrebbe gestire la richiesta basata sull'URL e il metodo HTTP.
    Utilizzo di wrap:
        L'handler selezionato è "avvolto" dalla funzione wrap. La funzione wrap può eseguire del codice prima di chiamare
		l'handler principale, come configurare il RequestContext o eseguire il logging.
        wrap può anche eseguire del codice dopo che l'handler principale ha finito, come il logging aggiuntivo o la gestione degli errori.
    Configurazione del RequestContext:
        All'interno della funzione wrap, un RequestContext viene creato e configurato. Questo potrebbe includere la
		generazione di un ID univoco per la richiesta e la configurazione di un logger per includere quell'ID nelle voci di log.
    Chiamata dell'Handler Principale:
        La funzione wrap chiama l'handler principale, passando la richiesta HTTP originale e il RequestContext configurato.
    Elaborazione della Richiesta:
        L'handler principale elabora la richiesta, utilizzando le informazioni nel RequestContext come necessario.
		Ad esempio, potrebbe utilizzare il logger nel RequestContext per registrare messaggi di log che includono l'ID della richiesta.
    Risposta:
        L'handler principale genera una risposta HTTP e la restituisce al client attraverso il Handler.
*/

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
// Handler restituisce un oggetto che implementa l'interfaccia http.Handler, che può essere utilizzato per gestire le richieste HTTP.

func (rt *_router) Handler() http.Handler {

	// Login enpoint
	rt.router.POST("/session", rt.wrap(rt.sessionHandler))

	// Search endpoint
	rt.router.GET("/users", rt.wrap(rt.getUsersQuery))

	// User Endpoint
	rt.router.PUT("/users/:id", rt.wrap(rt.putNickname))
	rt.router.GET("/users/:id", rt.wrap(rt.getUserProfile))

	// Ban endpoint
	rt.router.PUT("/users/:id/banned_users/:banned_id", rt.wrap(rt.putBan))
	rt.router.DELETE("/users/:id/banned_users/:banned_id", rt.wrap(rt.deleteBan))

	// Followers endpoint
	rt.router.PUT("/users/:id/followers/:follower_id", rt.wrap(rt.putFollow))
	rt.router.DELETE("/users/:id/followers/:follower_id", rt.wrap(rt.deleteFollow))

	// Stream endpoint
	rt.router.GET("/users/:id/home", rt.wrap(rt.getHome))

	// Photo Endpoint
	rt.router.POST("/users/:id/photos", rt.wrap(rt.postPhoto))
	rt.router.DELETE("/users/:id/photos/:photo_id", rt.wrap(rt.deletePhoto))
	rt.router.GET("/users/:id/photos/:photo_id", rt.wrap(rt.getPhoto))

	// Comments endpoint
	rt.router.POST("/users/:id/photos/:photo_id/comments", rt.wrap(rt.postComment))
	rt.router.DELETE("/users/:id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.deleteComment))

	// Likes endpoint
	rt.router.PUT("/users/:id/photos/:photo_id/likes/:like_id", rt.wrap(rt.putLike))
	rt.router.DELETE("/users/:id/photos/:photo_id/likes/:like_id", rt.wrap(rt.deleteLike))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
