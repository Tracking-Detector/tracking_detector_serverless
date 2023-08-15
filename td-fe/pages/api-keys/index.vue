<template>
    <div class="pa-4">
        <v-alert v-model="alert.isShowing" color="success" icon="$success" :title="alert.title" :text="alert.message"
            closable></v-alert>
        <h3 class="mb-4">API-KEYS</h3>
        <p style="width: 600px;" class="mb-4">On this page you can create and remove API-KEYS for users of the web
            extension.</p>
        <v-card>

            <v-list>
                <v-list-subheader>
                    Users
                </v-list-subheader>
                <v-divider></v-divider>
                <v-list-item v-for="user of userData">
                    <template v-slot:prepend>
                        <v-avatar :color="user.role == 'admin' ? 'red' : 'green'">
                            <v-icon color="white" v-if="user.role == 'admin'">mdi-security</v-icon>
                            <v-icon color="white" v-else>mdi-account</v-icon>
                        </v-avatar>
                    </template>
                    {{ user.email }}
                    <v-chip>
                        {{ user.role }}
                    </v-chip>
                    <template v-slot:append>
                        <v-btn v-if="user.role != 'admin'" color="red" icon="mdi-delete" variant="text"
                            @click="deleteUser(user['_id'])"></v-btn>
                    </template>
                </v-list-item>
            </v-list>
        </v-card>
        <v-card class="mt-2">
            <v-card-title>
                Create Users
            </v-card-title>
            <div class="ma-4">
                <v-text-field v-model="email" label="Email" required></v-text-field>
                <v-btn class="mt-2" @click="createNewUser">Create User</v-btn>
            </div>
        </v-card>
    </div>
</template>
<script setup>
const userData = ref([])
const email = ref("")
const isLoading = ref(true)
const config = useRuntimeConfig()
const alert = ref({
    isShowing: false,
    title: "User created",
    message: ""
})

const loadUserData = () => {
    isLoading.value = true
    fetch("/api/users", {
        headers: {
            "X-API-Key": 'Bearer ' + config.public.apiBase
        }
    }).then(response => {
        return response.json()
    }).then(body => {
        console.log(body)
        userData.value = body.data.data
        isLoading.value = false
    })
}

const createNewUser = () => {
    if (email.value == "") {
        alert.value.message = `No email provided`
        alert.value.isShowing = true
        return
    }
    isLoading.value = true
    fetch("/api/users", {
        method: "POST",
        headers: {
            "X-API-Key": 'Bearer ' + config.public.apiBase,
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            email: email.value
        })
    }).then(response => {

        return response.json()
    }).then(body => {
        loadUserData()
        alert.value.message = body.message + ` with key ${JSON.stringify(body.data)}`
        isLoading.value = false
        alert.value.isShowing = true
        email.value = ""
    })

}

const deleteUser = (id) => {
    fetch("/api/users/" + id, {
        method: "DELETE",
        headers: {
            "X-API-Key": 'Bearer ' + config.public.apiBase,
        },
    }).then(response => {

        return response.json()
    }).then(body => {
        loadUserData()
        console.log(JSON.stringify(body))
    })
}

onMounted(() => {

    loadUserData()
})
</script>