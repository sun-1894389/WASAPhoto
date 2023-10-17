<script>
// Rappresenta e gestisce la pagina di login
export default {
	data: function() {
		return {
			errormsg: null,
			identifier: "",
			disabled: true,
		}
	},
	methods: {
		async login() {
			// Effettua una richiesta POST all'endpoint "/session" con l'identificatore dell'utente.
			// this.loading = true;
			this.errormsg = null;
			try {
				// Login (POST): "/session"
				let response = await this.$axios.post("/session",{
					user_id: this.identifier.trim()
				});
				// Se la richiesta ha esito positivo, memorizza l'ID utente nel sessionStorage come 'token'.
				sessionStorage.setItem('token',response.data.user_id);
				// Reindirizza l'utente alla pagina "/home".
				this.$router.replace("/home")
				// informa che l'utente ha effettuato l'accesso
				this.$emit('updatedLoggedChild',true)
				
			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
		},
	},
	mounted(){
		if (sessionStorage.getItem('token')){
			this.$router.replace("/home")
		}
	},
	
}
</script>

<template>
	<div class="container-fluid h-100 m-0 p-0 login">

		<div class="row ">
			<div class="col">
				<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			</div>
		</div>

		<div class="row h-100 w-100 m-0">
			
			<form @submit.prevent="login" class="d-flex flex-column align-items-center justify-content-center p-0">

				<div class="row mt-2 mb-3 border-bottom">
					<div class="col">
						<h2 class="login-title">WASAPhoto Login</h2>
					</div>
				</div>

				<div class="row mt-2 mb-3">
					<div class="col">
						<input 
						type="text" 
						class="form-control" 
						v-model="identifier" 
						maxlength="16"
						minlength="3"
						placeholder="Your identifier" />
					</div>
				</div>

				<div class="row mt-2 mb-5 ">
					<div class="col ">
						<button class="btn btn-dark" :disabled="identifier == null || identifier.length >16 || identifier.length <3 || identifier.trim().length<3"> 
						Register/Login 
						</button>
					</div>
				</div>
			</form>
		</div>
	</div>
</template>

<style>
.login {
    background-image: url("../assets/images/BackgroundLogin.png");
    height: 100vh;
}

.login-title {
    color: black;
}
</style>
