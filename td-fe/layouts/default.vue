<template>
    <v-layout class="rounded rounded-md">
        <v-app-bar title="Tracking Detector Gym" color="red-darken-4"></v-app-bar>

        <v-navigation-drawer>
            <v-list>
                <v-list-item title="Overview" to="/">
                    <template v-slot:prepend>
                        <v-icon>mdi-view-dashboard</v-icon>
                    </template>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item title="Downloads" to="/downloads">
                    <template v-slot:prepend>
                        <v-icon>mdi-download-multiple</v-icon>
                    </template>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item title="Exports" to="/exports">
                    <template v-slot:prepend>
                        <v-icon>mdi-export</v-icon>
                    </template>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item title="Models" to="/models">
                    <template v-slot:prepend>
                        <v-icon>mdi-robot</v-icon>
                    </template>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item title="Users" to="/api-keys">
                    <template v-slot:prepend>
                        <v-icon>mdi-account-group</v-icon>
                    </template>
                </v-list-item>
                <v-divider></v-divider>
            </v-list>
            <template v-slot:append>
                <v-list>
                    <v-divider></v-divider>
                    <v-list-item title="Logs" href="/logs"></v-list-item>
                    <v-divider></v-divider>
                    <v-list-item title="Database" href="/mongo"></v-list-item>
                </v-list>
            </template>
        </v-navigation-drawer>

        <v-main class="ma-4">
            <slot />
        </v-main>
    </v-layout>
</template>