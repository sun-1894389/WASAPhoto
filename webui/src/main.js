// script che inizializza il frontend Vue e configura le sue parti principali

// importo metodi essenziali per creare un'applicazione Vue e per creare un oggetto reattivo.
import {
    createApp,
    reactive
} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Navbar from './components/Navbar.vue'
import Photo from './components/Photo.vue'
import UserMiniCard from './components/UserMiniCard.vue'
import PageNotFound from './components/PageNotFound.vue'
import LikeModal from './components/LikeModal.vue'
import CommentModal from './components/CommentModal.vue'
import PhotoComment from './components/PhotoComment.vue'

import './assets/dashboard.css'
import './assets/main.css'

// Crea un'istanza dell'applicazione Vue utilizzando il componente principale App.
const app = createApp(App)

// Aggiunge l'istanza di Axios all'applicazione come una proprietà globale,permettendo a qualsiasi componente di accedere ad Axios tramite this.$axios.
app.config.globalProperties.$axios = axios;

// registro i vari component
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Navbar", Navbar);
app.component("Photo", Photo);
app.component("UserMiniCard", UserMiniCard);
app.component("PageNotFound", PageNotFound);
app.component("LikeModal", LikeModal);
app.component("CommentModal", CommentModal);
app.component("PhotoComment", PhotoComment);

// Aggiunge il router all'applicazione, permettendo la navigazione tra le diverse viste.
app.use(router)

// Monta l'applicazione Vue sull'elemento con l'ID app
app.mount('#app')