<template>
  <div>
    <h3 class="mb-4 text-h5">Welcome in the Tracking Detector Gym</h3>
    <v-card class="pa-4 mt-3">
      <v-card-title>Health</v-card-title>
      <v-divider></v-divider>
      <v-list>
        <v-list-item
          v-for="check of performedHealthChecks"
          :key="check.service"
        >
          <template v-slot:prepend>
            <v-icon v-if="check.status == 200" color="green"> mdi-check</v-icon>
            <v-icon v-else color="red"> mdi-alert-circle</v-icon>
          </template>
          <v-list-item-title>{{ check.service }}</v-list-item-title>
          <v-list-item-subtitle>{{ check.message }}</v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card>
    <v-card class="pa-4 mt-3">
      <v-card-title>API</v-card-title>
      <v-divider></v-divider>
      <v-card-title class="mt-3 mb-2">Requests</v-card-title>
      <v-table>
        <thead>
          <tr>
            <th class="text-left">Endpoint</th>
            <th class="text-left">Method</th>
            <th class="text-left">Description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in requestsEndpoint" :key="item.description">
            <td>
              <a :href="item.url">{{ item.url }}</a>
            </td>
            <td>{{ item.method }}</td>
            <td>{{ item.description }}</td>
          </tr>
        </tbody>
      </v-table>
      <v-card-title class="mt-3 mb-2">Users</v-card-title>
      <v-table>
        <thead>
          <tr>
            <th class="text-left">Endpoint</th>
            <th class="text-left">Method</th>
            <th class="text-left">Description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in usersEndpoint" :key="item.description">
            <td>
              <a :href="item.url">{{ item.url }}</a>
            </td>
            <td>{{ item.method }}</td>
            <td>{{ item.description }}</td>
          </tr>
        </tbody>
      </v-table>
      <v-card-title class="mt-3 mb-2">Dispatch</v-card-title>
      <v-table>
        <thead>
          <tr>
            <th class="text-left">Endpoint</th>
            <th class="text-left">Method</th>
            <th class="text-left">Description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in dispatchEndpoint" :key="item.description">
            <td>
              <a :href="item.url">{{ item.url }}</a>
            </td>
            <td>{{ item.method }}</td>
            <td>{{ item.description }}</td>
          </tr>
        </tbody>
      </v-table>
      <v-card-title class="mt-3 mb-2">Download</v-card-title>
      <v-table>
        <thead>
          <tr>
            <th class="text-left">Endpoint</th>
            <th class="text-left">Method</th>
            <th class="text-left">Description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in downloadEndpoint" :key="item.description">
            <td>
              <a :href="item.url">{{ item.url }}</a>
            </td>
            <td>{{ item.method }}</td>
            <td>{{ item.description }}</td>
          </tr>
        </tbody>
      </v-table>
      <v-card-title class="mt-3 mb-2">Training Runs</v-card-title>
      <v-table>
        <thead>
          <tr>
            <th class="text-left">Endpoint</th>
            <th class="text-left">Method</th>
            <th class="text-left">Description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in trainingRunsEndpoint" :key="item.description">
            <td>
              <a :href="item.url">{{ item.url }}</a>
            </td>
            <td>{{ item.method }}</td>
            <td>{{ item.description }}</td>
          </tr>
        </tbody>
      </v-table>
    </v-card>
  </div>
</template>
<script setup>
useMeta({
  title: "Tracking Detector",
});
const config = useRuntimeConfig();
const healthChecks = ref([
  {
    service: "Users",
    endpoint: "/api/users/health",
  },
  {
    service: "Requests",
    endpoint: "/api/requests/health",
  },
  {
    service: "Dispatch",
    endpoint: "/api/dispatch/health",
  },
  {
    service: "Download",
    endpoint: "/api/transfer/health",
  },
  {
    service: "Training-Runs",
    endpoint: "/api/training-runs/health",
  },
]);

const performedHealthChecks = ref([]);

onMounted(() => {
  healthChecks.value.map((check) => {
    fetch(check.endpoint, {
      headers: {
        "X-API-Key": "Bearer " + config.public.apiBase,
      },
    })
      .then((response) => {
        if (response.status != 200) {
          return { status: 500, message: "Service unavailable." };
        }
        return response.json();
      })
      .then((body) => {
        performedHealthChecks.value.push({
          service: check.service,
          status: body.status,
          message: body.message,
        });
      });
  });
});

const usersEndpoint = ref([
  {
    url: "/api/users/health",
    method: "GET",
    description: "Endpoint to check the health of the users microservice.",
  },
  {
    url: "/api/users",
    method: "GET",
    description: "Get all users which are registered in the application.",
  },
  {
    url: "/api/users",
    method: "POST",
    description: "Create a new client that can send request data into the db.",
  },
  {
    url: "/api/users/:userId",
    method: "DELETE",
    description: "Delete a user by id.",
  },
]);
const requestsEndpoint = ref([
  {
    url: "/api/requests/health",
    method: "GET",
    description: "Endpoint to check the health of the requests microservice.",
  },
  {
    url: "/api/requests/:requestId",
    method: "GET",
    description: "Get Request data by id in JSON format.",
  },
  {
    url: "/api/requests",
    method: "POST",
    description: "Create a single request object.",
  },
  {
    url: "/api/requests/multiple",
    method: "POST",
    description: "Create multiple request objects. Ideal to seed the db.",
  },
  {
    url: "/api/requests?url=url&page=1&pageSize=10",
    method: "GET",
    description: "Endpoint to search for request data inside the db.",
  },
]);
const dispatchEndpoint = ref([
  {
    url: "/api/dispatch/health",
    method: "GET",
    description: "Endpoint to check the health of the dispatch microservice.",
  },
  {
    url: "/api/dispatch/export",
    method: "GET",
    description: "Returns the available exports.",
  },
  {
    url: "/api/dispatch/model",
    method: "GET",
    description: "Returns the available models.",
  },
  {
    url: "/api/dispatch/export/:extractorName",
    method: "POST",
    description: "Dispatches an export job to the redis queue.",
  },
  {
    url: "/api/dispatch/train/:modelName/run/:dataSetName",
    method: "POST",
    description: "Dispatches a training job to the redis queue.",
  },
]);
const downloadEndpoint = ref([
  {
    url: "/api/transfer/health",
    method: "GET",
    description: "Endpoint to check the health of the transfer microservice.",
  },
  {
    url: "/api/transfer/export/:fileName",
    method: "GET",
    description: "Downloads a certain export as .tar.gz.",
  },
  {
    url: "/api/transfer/export",
    method: "GET",
    description: "Returns all the available export downloads.",
  },
  {
    url: "/api/transfer/models/:modelName/:zippedModelName",
    method: "GET",
    description: "Downloads a certain model as .tar.gz.",
  },
  {
    url: "/api/transfer/models/:modelName/:trainingSet/:filename",
    method: "GET",
    description:
      "Downloads the raw files from the model. This can be used to load the model via tensorflow.",
  },
  {
    url: "/api/transfer/models",
    method: "GET",
    description: "Returns all the available model downloads.",
  },
]);
const trainingRunsEndpoint = ref([
  {
    url: "/api/training-runs/health",
    method: "GET",
    description:
      "Endpoint to check the health of the training-runs microservice.",
  },
  {
    url: "/api/training-runs",
    method: "GET",
    description: "Returns all training runs as a list.",
  },
  {
    url: "/api/training-runs/:modelName",
    method: "GET",
    description: "Returns all training-runs for a specific model.",
  },
]);
</script>
