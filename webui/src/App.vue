<script setup>
// layout principale dell'applicazione
// importo Routerlink,Routerview per la gestione delle vie. 
// RouterView è un contenitore per il componente corrispondente alla via corrente.
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data(){
		return{
			// indica se l'utente è attualmente autenticato.
			logged: false,
			searchValue: "",
		}
	},
	methods:{
		logout(newValue){
			// Imposta logged su false e reindirizza l'utente alla pagina di login.
			this.logged = newValue
			this.$router.replace("/login")
		},
		updateLogged(newLogged){
			// Aggiorna lo stato logged con il nuovo valore.
			this.logged = newLogged
		},
		updateView(newRoute){
			// Reindirizza l'utente alla nuova via specificata.
			this.$router.replace(newRoute)
		},
		search(queryParam){
			// mposta il valore di searchValue e reindirizza l'utente alla pagina di ricerca.
			this.searchValue= queryParam
			this.$router.replace("/search")
		},
	},

	
	created(){
		// Quando il componente viene creato, controlla se esiste un elemento notFirstStart
		// nel sessionStorage. Se non esiste, pulisce il sessionStorage e imposta notFirstStart su true.
		if (!sessionStorage.getItem('notFirstStart')){
			sessionStorage.clear()
			sessionStorage.setItem('notFirstStart',true)
			// console.log("first start")
		}
		
	},
	

	mounted(){

		// controlla se esiste un token nel sessionStorage
		// Se esiste, imposta logged su true, altrimenti reindirizza l'utente alla pagina di login.
		if (!sessionStorage.getItem('token')){
			this.$router.replace("/login")
		}else{
			this.logged = true
		}
	},
}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<div class="col p-0">
				<main >
					<Navbar v-if="logged" 
					@logoutNavbar="logout" 
					@requestUpdateView="updateView"
					@searchNavbar="search"/>

					<RouterView 
					@updatedLoggedChild="updateLogged" 
					@requestUpdateView="updateView"
					:searchValue="searchValue"/>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>
