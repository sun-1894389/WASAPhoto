<script>
// view per ricerca per gli utenti 
export default {
	data: function() {
		return {
			// lista di utenti che corrispondono ai criteri di ricerca.
			users: [],
			errormsg: null,
		}
	},

	// Il valore di ricerca inserito dall'utente.
	props:['searchValue'],

	watch:{
		// Ogni volta che il valore di ricerca cambia, viene chiamata la funzione 
		searchValue: function(){
			this.loadSearchedUsers()
		},
	},

	methods:{
		//  Carica gli utenti che corrispondono ai criteri di ricerca.
		async loadSearchedUsers(){
			this.errormsg = null;
			if (
				this.searchValue === undefined ||
				this.searchValue === "" || 
				this.searchValue.includes("?") ||
				this.searchValue.includes("_")){
				this.users = []
				return 
			}
			try {
				// Search user (PUT):  "/users"
				let response = await this.$axios.get("/users",{
						params: {
						id: this.searchValue,
					},
				});
				this.users = response.data

			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		// Reindirizza l'utente al profilo dell'utente selezionato.
		goToProfile(profileId){
			this.$router.replace("/users/"+profileId)
		}
	},

	async mounted(){
		// Check if the user is logged
		if (!sessionStorage.getItem('token')){
			this.$router.replace("/login")
		}
		await this.loadSearchedUsers()
		
	},
}
</script>

<template>
	<div class="container-fluid h-100 ">
		<UserMiniCard v-for="(user,index) in users" 
		:key="index"
		:identifier="user.user_id" 
		:nickname="user.nickname" 
		@clickedUser="goToProfile"/>

		<p v-if="users.length == 0" class="no-result-text d-flex justify-content-center"> No users found.</p>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>

.no-result-text{
	color: white;
	font-style: italic;
}
</style>
