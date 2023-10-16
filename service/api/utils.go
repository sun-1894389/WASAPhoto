package api

import (
	"net/http"
	"strings"
)

// funzione che verifica se l'identificatore di un utente ha la lunghezza corretta.
func validIdentifier(identifier string) bool {
	// Rimuovo degli spazi bianchi all'inizio e alla fine dell'identificatore.
	// Controllo se la lunghezza dell'identificatore è compresa tra 3 e 16 caratteri,se non è vuoto e se non contiene i caratteri "?" o "_".
	var trimmedId = strings.TrimSpace(identifier)
	return len(identifier) >= 3 && len(identifier) <= 16 && trimmedId != "" && !strings.ContainsAny(trimmedId, "?_")
}

// funzione che estrae il token Bearer dall'intestazione di autorizzazione.
func extractBearer(authorization string) string {
	// Divido l'intestazione di autorizzazione in token utilizzando lo spazio come delimitatore.
	var tokens = strings.Split(authorization, " ")
	// Se ci sono esattamente due token, restituisci il secondo token (il token Bearer) dopo aver rimosso eventuali spazi.
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	// Se non ci sono 2 token restituisco stringa vuota
	return ""
}

// Funzione che verifica se l'utente che effettua la richiesta ha un token valido per l'endpoint specificato.Restituisce 0 se è valido,o errore
func validateRequestingUser(identifier string, bearerToken string) int {

	// Se l'utente che effettua la richiesta ha un token non valido, restituisci un codice di stato HTTP 403
	if isNotLogged(bearerToken) {
		return http.StatusForbidden
	}

	// Se l'ID dell'utente che effettua la richiesta è diverso da quello nel percorso, restituisci un codice di stato HTTP 401
	if identifier != bearerToken {
		return http.StatusUnauthorized
	}
	return 0
}

// funzione che verifica se un utente è loggato.
// Restituisci true se la stringa di autenticazione è vuota (cioè l'utente non è loggato), altrimenti restituisci false
func isNotLogged(auth string) bool {

	return auth == ""
}
