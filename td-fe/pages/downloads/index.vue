<template>
  <div>
    <h3 class="mb-4 text-h5">Welcome on Downloads</h3>
    <p style="width: 600px" class="mb-4 text-body-1">
      Here are all the available downloadable files in one place. There a
      dataset exports available in the export directory and trained models in
      the models directory.
    </p>

    <v-card elevation="2">
      <v-toolbar class="pl-4"> files://{{ dir }} </v-toolbar>
      <v-row>
        <v-col cols="4" style="border-right: 1px solid lightgray">
          <v-list>
            <div v-for="item in folderStructure" :key="item.name">
              <v-list-group
                v-if="item.children != undefined"
                :value="item.name"
              >
                <template v-slot:activator="{ props }">
                  <v-list-item
                    v-bind="props"
                    prepend-icon="mdi-folder"
                    :title="item.name"
                    @click="changeDir(item)"
                  ></v-list-item>
                </template>
                <div v-for="subItem in item.children" :key="subItem.name">
                  <v-list-group
                    v-if="subItem.children != undefined"
                    :value="subItem.name"
                  >
                    <template v-slot:activator="{ props }">
                      <v-list-item
                        v-bind="props"
                        prepend-icon="mdi-folder"
                        :title="subItem.name"
                        @click="changeDir(subItem)"
                      ></v-list-item>
                    </template>
                    <div
                      v-for="subSubItem in subItem.children"
                      :key="subSubItem.name"
                    >
                      <v-list-group
                        v-if="subSubItem.children != undefined"
                        :value="subSubItem.name"
                      >
                        <template v-slot:activator="{ props }">
                          <v-list-item
                            v-bind="props"
                            prepend-icon="mdi-folder"
                            :title="subSubItem.name"
                            @click="changeDir(subSubItem)"
                          ></v-list-item>
                        </template>
                        <v-list-item
                          v-for="subSubSubItem in subSubItem.children"
                          :key="subSubSubItem.name"
                          @click="changeDir(subSubSubItem)"
                        >
                          <template v-slot:prepend>
                            <v-icon>mdi-file</v-icon>
                          </template>
                          {{ subSubSubItem.name }}
                        </v-list-item>
                      </v-list-group>
                      <v-list-item v-else @click="changeDir(subSubItem)">
                        <template v-slot:prepend>
                          <v-icon>mdi-file</v-icon>
                        </template>
                        {{ subSubItem.name }}
                      </v-list-item>
                    </div>
                  </v-list-group>
                  <v-list-item v-else @click="changeDir(subItem)">
                    <template v-slot:prepend>
                      <v-icon>mdi-file</v-icon>
                    </template>
                    {{ subItem.name }}
                  </v-list-item>
                </div>
              </v-list-group>
              <v-list-item v-else @click="changeDir(item)">
                {{ item.name }}
              </v-list-item>
              <v-divider></v-divider>
            </div>
          </v-list>
        </v-col>
        <v-col cols="8">
          <v-list>
            <div v-for="entry in currentDir.children" :key="entry.name">
              <v-list-item
                v-if="entry.children != undefined"
                @click="changeDir(entry)"
              >
                <template v-slot:prepend>
                  <v-icon>mdi-folder</v-icon>
                </template>
                {{ entry.name }}
              </v-list-item>
              <v-list-item v-else>
                <template v-slot:prepend>
                  <v-icon>mdi-file</v-icon>
                </template>
                {{ entry.name }}
                <template v-slot:append>
                  <v-btn
                    icon="mdi-download"
                    variant="text"
                    :href="'/api/transfer/' + entry.path"
                    download
                  ></v-btn>
                </template>
              </v-list-item>
              <v-divider></v-divider>
            </div>
          </v-list>
        </v-col>
      </v-row>
    </v-card>
  </div>
</template>
<script setup>
useHead({
  title: "Tracking Detector - Downloads",
});
const downloadableData = ref({});
const folderStructure = ref([]);
const currentDir = ref([]);
const config = useRuntimeConfig();
const breadCrumps = ref([
  {
    title: "Files",
    disabled: false,
    href: "breadcrumbs_dashboard",
  },
]);

const dir = ref("");

const changeDir = (item) => {
  if (item.children != undefined) {
    dir.value = item.path;
  } else {
    const path = item.path.split("/");
    dir.value = path.slice(0, path.length - 1).join("/");
  }
};

const transformToFolderArchitecture = () => {
  const getSubfolder = (data, key, path) => {
    const keys = Object.keys(data[key]);
    if (keys.length > 0) {
      const children = keys.map((x) =>
        getSubfolder(data[key], x, path + "/" + x)
      );
      return {
        name: key,
        children: children,
        path: path,
      };
    }
    const ext = key.split(".").filter(Boolean).slice(1).join(".");
    return {
      name: key,
      file: ext,
      path: path,
    };
  };
  folderStructure.value = Object.keys(downloadableData.value).map((x) =>
    getSubfolder(downloadableData.value, x, x)
  );
};
onMounted(() => {
  fetch("/api/transfer/models", {
    headers: {
      "X-API-Key": "Bearer " + config.public.apiBase,
    },
  })
    .then((response) => {
      if (response.status != 200) {
        return undefined;
      } else {
        return response.json();
      }
    })
    .then((body) => {
      if (body != undefined) {
        downloadableData.value.models = body;
      }

      transformToFolderArchitecture();
    });
  fetch("/api/transfer/export", {
    headers: {
      "X-API-Key": "Bearer " + config.public.apiBase,
    },
  })
    .then((response) => {
      if (response.status != 200) {
        return undefined;
      } else {
        return response.json();
      }
    })
    .then((body) => {
      if (body != undefined) {
        downloadableData.value.export = body;
      }
      transformToFolderArchitecture();
    });
});
watch(dir, (newVal) => {
  const keys = newVal.split("/");
  let curDir = folderStructure.value;
  for (let i = 0; i < keys.length; i++) {
    if (i == 0) {
      curDir = curDir.find((x) => x.name == keys[i]);
    } else {
      curDir = curDir.children.find((x) => x.name == keys[i]);
    }
  }
  currentDir.value = curDir;
});
</script>
