<template>
    <v-layout class="rounded rounded-md">
        <v-app-bar title="Tracking Detector Gym" color="red-darken-4"></v-app-bar>

        <v-navigation-drawer>
            <v-list>
                <v-list-item title="Overview" to="/"></v-list-item>
                <v-divider></v-divider>
                <v-list-item title="Downloads" to="/downloads"></v-list-item>
                <v-divider></v-divider>
                <v-list-item title="Exports" to="/exports"></v-list-item>
                <v-divider></v-divider>
                <v-list-item title="Models" to="/models"></v-list-item>
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