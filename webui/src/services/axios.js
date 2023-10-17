// axios.js crea e configura un'istanza di axios,per effettuare richieste HTTP.
// se un token è presente nel sessionStorage, lo aggiunge automaticamente all'header Authorization di ogni richiesta in uscita.
import axios from "axios";

//Crea un'istanza di axios
const instance = axios.create({
    baseURL: __API_URL__,
    timeout: 1000 * 5,
});

// Aggiunge un "interceptor" alle richieste in uscita. Questo permette di modificare
// le richieste o le risposte prima che vengano inviate o dopo che sono state ricevute.
instance.interceptors.request.use(
    // gestisco le richiste in uscita, recupero il token di autenticazione/uuid da session storage
    (config) => {
        const token = sessionStorage.getItem('token');
        // Se il token esiste, lo aggiunge all'header Authorization delle richieste in uscita con il prefisso 'Bearer'.
        if (token) {
            config.headers['Authorization'] = 'Bearer ' + token;
        }

        return config
    },
    // gestisco l'errore
    (error) => {
        return Promise.reject(error);
    }
)

export default instance;