<template>
  <el-row class="tac">
    <el-col :lg="4" :span="6" class="border-right">
      <el-menu
          class="no-border-right"
          default-active="allHosts"
          @close="handleMenuSelect"
          @open="handleMenuSelect"
          @select="handleMenuSelect"
      >
        <el-menu-item index="allHosts">
          <el-icon>
            <document/>
          </el-icon>
          <span>All Hosts</span>
        </el-menu-item>
        <el-menu-item index="inUse">
          <el-icon>
            <el-icon-document-checked/>
          </el-icon>
          <span>In Use</span>
        </el-menu-item>
        <el-sub-menu index="hostGroups">
          <template #title>
            <el-icon>
              <icon-menu/>
            </el-icon>
            <span>Host Groups</span>
          </template>
          <el-menu-item v-for="(group, groupName) in listByGroup.list" :index="'showGroup:'+groupName"
                        @mouseleave="group.mouseover=false" @mouseover="group.mouseover=true">
            <template #title>
              <el-checkbox v-model="group.show" :indeterminate="group.showNum>0 && group.hideNum>0"
                           :checked="group.showNum>0 && group.hideNum===0" @change="handleSwitchByGroupName(group)"
                           @click.stop class="mr-1"/>
              {{ groupName }}
              <Edit v-if="group.mouseover" style="width: 1em; height: 1em; margin-left: 8px"
                    @click="openDialogForm(groupName)" @click.stop/>
            </template>
          </el-menu-item>
        </el-sub-menu>
        <el-menu-item index="addHost">
          <el-icon>
            <document-add/>
          </el-icon>
          <span>Add Host</span>
        </el-menu-item>
      </el-menu>
    </el-col>
    <el-col :span="18" class="show-content">
      <el-tabs v-model="activeTabName" class="ml-2 mr-2" @tab-click="handleClickTab">
        <el-tab-pane label="Normal" name="normal">
          <el-table v-if="activeMenuIndex==='allHosts'" :key="groupName"
                    :data="filterTableData"
                    :model="allHostsForm"
                    @selection-change="handleSelectionChange"
                    stripe
                    style="width: 100%"
          >
            <el-table-column type="selection" width="30"/>
            <el-table-column align="center" label="Active" width="80">
              <template #default="scope">
                <el-switch v-model="scope.row.show" @change="handleSwitchByHostnameId(scope.row)"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="IP" prop="ip" width="150">
              <template #header>
                <el-input v-model="filterIP" placeholder="IP" size="small"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="Hostname" prop="hostName">
              <template #header>
                <el-input v-model="filterHostname" placeholder="Hostname" size="small"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="Group Name" prop="groupName">
              <template #header>
                <el-input v-model="filterGroupName" placeholder="Group Name" size="small"/>
              </template>
            </el-table-column>
            <el-table-column label="Action" width="180">
              <template #header>
                <el-dropdown>
                <span class="el-dropdown-link">
                  Action
                  <el-icon class="el-icon--right">
                    <arrow-down/>
                  </el-icon>
                </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item>
                        <el-link type="primary" @click="handleDeleteAllSelected()" @click.stop>Delete All Selected
                        </el-link>
                      </el-dropdown-item>
                      <el-dropdown-item>
                        <el-link type="primary" @click="dialogMoveSelectedFormVisible=true" @click.stop>Move All
                          Selected
                        </el-link>
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
              <template #default="scope">
                <el-button size="small" type="danger" @click="handleDeleteHost(scope.$index, scope.row)">Delete
                </el-button>
                <el-button size="small" type="primary" @click="openDialogMoveForm(scope.$index, scope.row)">Move
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-form v-else-if="activeMenuIndex==='addHost'" :model="addHostForm" label-width="120px">
            <el-form-item label="Group Name">
              <el-autocomplete
                  v-model="addHostForm.groupName"
                  :fetch-suggestions="querySearchGroupNames"
                  class="el-col-10"
                  clearable
                  placeholder="Group Name"
              />
            </el-form-item>
            <el-form-item label="IP">
              <el-input v-model="addHostForm.ip" class="el-col-10" placeholder="e.g. 127.0.0.1"/>
            </el-form-item>
            <el-form-item label="Host">
              <el-input v-model="addHostForm.hostName" class="el-col-10" placeholder="e.g. www.domain.com"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onSubmitAddHost">Add</el-button>
            </el-form-item>
          </el-form>
          <el-table v-else-if="activeMenuIndex==='showGroup'"
                    :key="groupName"
                    :data="filterTableData"
                    @selection-change="handleSelectionChange"
                    stripe
                    style="width: 100%"
          >
            <el-table-column type="selection" width="30"/>
            <el-table-column align="center" label="Active" width="80">
              <template #default="scope">
                <el-switch v-model="scope.row.show" @change="handleSwitchByHostnameId(scope.row)"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="IP" prop="ip" width="150">
              <template #header>
                <el-input v-model="filterIP" placeholder="IP" size="small"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="Hostname" prop="hostName">
              <template #header>
                <el-input v-model="filterHostname" placeholder="Hostname" size="small"/>
              </template>
            </el-table-column>
            <el-table-column label="Action" width="180">
              <template #header>
                <el-dropdown>
                <span class="el-dropdown-link">
                  Action
                  <el-icon class="el-icon--right">
                    <arrow-down/>
                  </el-icon>
                </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item>
                        <el-link type="primary" @click="handleDeleteAllSelected()" @click.stop>Delete All Selected
                        </el-link>
                      </el-dropdown-item>
                      <el-dropdown-item>
                        <el-link type="primary" @click="dialogMoveSelectedFormVisible=true" @click.stop>Move All
                          Selected
                        </el-link>
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
              <template #default="scope">
                <!--            <el-button size="small" @click="handleEditHost(scope.$index, scope.row)">Edit</el-button>-->
                <el-button size="small" type="danger" @click="handleDeleteHost(scope.$index, scope.row)">Delete
                </el-button>
                <el-button size="small" type="primary" @click="openDialogMoveForm(scope.$index, scope.row)">Move
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-table v-else-if="activeMenuIndex==='inUse'"
                    :key="groupName"
                    :data="filterTableData"
                    @selection-change="handleSelectionChange"
                    stripe
                    style="width: 100%"
          >
            <el-table-column type="selection" width="30"/>
            <el-table-column align="center" label="Active" width="80">
              <template #default="scope">
                <el-switch v-model="scope.row.show" @change="handleSwitchByHostnameId(scope.row)"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="IP" prop="ip" width="150">
              <template #header>
                <el-input v-model="filterIP" placeholder="IP" size="small"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="Hostname" prop="hostName">
              <template #header>
                <el-input v-model="filterHostname" placeholder="Hostname" size="small"/>
              </template>
            </el-table-column>
            <el-table-column align="center" label="Group Name" prop="groupName">
              <template #header>
                <el-input v-model="filterGroupName" placeholder="Group Name" size="small"/>
              </template>
            </el-table-column>
            <el-table-column label="Action" width="180">
              <template #header>
                <el-dropdown>
                <span class="el-dropdown-link">
                  Action
                  <el-icon class="el-icon--right">
                    <arrow-down/>
                  </el-icon>
                </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item>
                        <el-link type="primary" @click="handleDeleteAllSelected()" @click.stop>Delete All Selected
                        </el-link>
                      </el-dropdown-item>
                      <el-dropdown-item>
                        <el-link type="primary" @click="dialogMoveSelectedFormVisible=true" @click.stop>Move All
                          Selected
                        </el-link>
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
              <template #default="scope">
                <!--            <el-button size="small" @click="handleEditHost(scope.$index, scope.row)">Edit</el-button>-->
                <el-button size="small" type="danger" @click="handleDeleteHost(scope.$index, scope.row)">Delete
                </el-button>
                <el-button size="small" type="primary" @click="openDialogMoveForm(scope.$index, scope.row)">Move
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="Advanced" name="advanced">
          <el-form v-if="activeMenuIndex==='allHosts'" :model="allHostsForm">
            <el-form-item>
              <el-alert class="el-form-item__content" effect="dark"
                        title="e.g. 127.0.0.1 www.domain.com # Group Name One # Group Name Two" type="info"/>
            </el-form-item>
            <el-form-item class="mt-2">
              <CodeEditor v-model="allHostsForm.text"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onSubmitAllHosts">Save</el-button>
              <el-button @click="copyToClipboard(allHostsForm)">Copy to Clipboard</el-button>
            </el-form-item>
          </el-form>
          <el-form v-else-if="activeMenuIndex==='addHost'" :model="addHostsTextForm">
            <el-form-item>
              <el-alert class="el-form-item__content" effect="dark"
                        title="e.g. 127.0.0.1 www.domain.com # Group Name One # Group Name Two" type="info"/>
            </el-form-item>
            <el-form-item class="mt-2">
              <CodeEditor v-model="addHostsTextForm.text"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onSubmitAddHostsTextForm">Save</el-button>
            </el-form-item>
          </el-form>
          <el-form v-else-if="activeMenuIndex==='showGroup'" :model="allGroupHostsForm">
            <el-form-item>
              <el-alert class="el-form-item__content" effect="dark"
                        title="e.g. 127.0.0.1 www.domain.com # Group Name One # Group Name Two" type="info"/>
            </el-form-item>
            <el-form-item class="mt-2">
              <CodeEditor v-model="allGroupHostsForm.text"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onSubmitAllGroupHosts">Save</el-button>
              <el-button @click="copyToClipboard(allGroupHostsForm)">Copy to Clipboard</el-button>
            </el-form-item>
          </el-form>
          <el-form v-else-if="activeMenuIndex==='inUse'" :model="allInUseHostsForm">
            <el-form-item>
              <el-alert class="el-form-item__content" effect="dark"
                        title="e.g. 127.0.0.1 www.domain.com # Group Name One # Group Name Two" type="info"/>
            </el-form-item>
            <el-form-item class="mt-2">
              <CodeEditor v-model="allInUseHostsForm.text"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onSubmitAllInUseHosts">Save</el-button>
              <el-button @click="copyToClipboard(allInUseHostsForm)">Copy to Clipboard</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-col>
  </el-row>
  <el-dialog v-model="dialogFormVisible" title="Change group name">
    <el-form :model="dialogForm">
      <el-form-item label="Group Name">
        <el-autocomplete
            v-model="dialogForm.groupName"
            :fetch-suggestions="querySearchGroupNames"
            clearable
            placeholder="Group Name"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitDialogForm">
          Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>
  <el-dialog v-model="dialogMoveFormVisible" title="Move to">
    <el-form :model="dialogForm">
      <el-form-item label="Group Name">
        <el-autocomplete
            v-model="dialogForm.groupName"
            :fetch-suggestions="querySearchGroupNames"
            clearable
            placeholder="Group Name"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogMoveFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitDialogMoveForm">
          Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>
  <el-dialog v-model="dialogMoveSelectedFormVisible" title="Move all selected to">
    <el-form :model="dialogForm">
      <el-form-item label="Group Name">
        <el-autocomplete
            v-model="dialogForm.groupName"
            :fetch-suggestions="querySearchGroupNames"
            clearable
            placeholder="Group Name"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogMoveSelectedFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitDialogMoveSelectedForm">
          Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {ArrowDown, Document, DocumentAdd, Edit, Menu as IconMenu} from '@element-plus/icons-vue'
import CodeEditor from './components/CodeEditor.vue'
import {computed, onMounted, reactive, ref} from "vue"
import type {TabsPaneContext} from 'element-plus'
import {ElMessage, ElMessageBox, ElTable} from 'element-plus'
import 'element-plus/dist/index.css'
import {
  AddHost,
  DeleteHost,
  DeleteHosts,
  GetAllGroupNames,
  GetHostsText,
  GetInUseHostsText,
  GetMyHosts,
  SaveAddHostsText,
  SaveAllGroupHosts,
  SaveAllHosts,
  SaveAllInUseHosts,
  SetGroupName,
  SetGroupNameByHostnameId,
  SetGroupNameByHosts,
  SwitchByGroupName,
  SwitchByHostnameId
} from "../wailsjs/go/main/App"
import {ClipboardSetText} from "../wailsjs/runtime"

const filterIP = ref('')
const filterHostname = ref('')
const filterGroupName = ref('')
const activeTabName = ref('normal')
const groupName = ref('')
const activeMenuIndex = ref('allHosts')
const listByGroup = reactive({list: {}, hostsText: ''})
const addHostForm = reactive({
  groupName: '',
  ip: '',
  hostName: '',
})
const addHostsTextForm = reactive({
  text: ''
})
const allHostsForm = reactive({
  text: ''
})
const allInUseHostsForm = reactive({
  text: ''
})
const allGroupHostsForm = reactive({
  text: ''
})
const tableData = ref([])
const allGroupNames = ref([])
const dialogFormVisible = ref(false)
const dialogMoveFormVisible = ref(false)
const dialogMoveSelectedFormVisible = ref(false)
const dialogForm = reactive({
  index: 0,
  id: 0,
  oldGroupName: '',
  groupName: '',
})

function getHostsText() {
  GetHostsText().then(result => {
    allHostsForm.text = result
  })
}

function getInUseHostsText() {
  GetInUseHostsText().then(result => {
    allInUseHostsForm.text = result
  })
}

function getListByGroup() {
  GetMyHosts().then(result => {
    console.log("listByGroup", result.listByGroup)
    listByGroup.list = result.listByGroup

    if (groupName.value !== "" && activeMenuIndex.value === "showGroup") {
      tableData.value.splice(0, tableData.value.length)
      let groupInfo = listByGroup.list[groupName.value].list
      for (let k in groupInfo) {
        tableData.value.push(groupInfo[k])
      }
      allGroupHostsForm.text = listByGroup.list[groupName.value].hostsText
      console.log(tableData)
    }
  })
}

function getAllList() {
  GetMyHosts().then(result => {
    console.log("list", result.list)

    tableData.value.splice(0, tableData.value.length)
    for (let k in result.list) {
      tableData.value.push(result.list[k])
    }
  })
}

function getActiveList() {
  GetMyHosts().then(result => {
    console.log("list", result.list)

    tableData.value.splice(0, tableData.value.length)
    for (let k in result.list) {
      if (result.list[k].show) {
        tableData.value.push(result.list[k])
      }
    }
  })
}

const handleMenuSelect = (key: string, keyPath: string[]) => {
  console.log('select', key, keyPath)
  if (key === 'allHosts') {
    activeMenuIndex.value = key
    if (activeTabName.value === "normal") {
      getAllList()
    } else {
      getHostsText()
    }
  } else if (key === 'addHost') {
    activeMenuIndex.value = key
    getAllGroupNames()
  } else if (key === 'inUse') {
    activeMenuIndex.value = key
    groupName.value = key
    if (activeTabName.value === "normal") {
      getActiveList()
    } else {
      getInUseHostsText()
    }
  } else if (key.substring(0, 9) === 'showGroup') {
    activeMenuIndex.value = key.substring(0, 9)
    tableData.value.splice(0, tableData.value.length)
    if (key.substring(10) !== "") {
      groupName.value = key.substring(10)
    }
  }
  getListByGroup()
}

const handleSwitchByGroupName = (group) => {
  console.log('switch', group.groupName, group.show)
  SwitchByGroupName(group.groupName, group.show).then(result => {
    if (result !== '') {
      ElMessage.error('switch failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('switch successfully!')
    }
  })
}

const handleSwitchByHostnameId = (row) => {
  console.log('switch', row.groupName, row.id, row.hostName, row.show)
  SwitchByHostnameId(row.groupName, row.id, row.show).then(result => {
    if (result !== '') {
      ElMessage.error('switch failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('switch successfully!')
    }
  })
}

const handleDeleteHost = (index, row) => {
  console.log('delete', index, row.groupName, row.hostName, row.id)
  ElMessageBox.confirm('Are you sure to delete this host?')
      .then(() => {
        DeleteHost(row.groupName, row.id).then(result => {
          if (result !== '') {
            ElMessage.error('delete failed!' + result)
          } else {
            tableData.value.splice(index, 1)
            handleMenuSelect(activeMenuIndex.value, null)
            ElMessage.success('delete successfully!')
          }
        })
      })
      .catch(() => {
        // catch error
      })
}

const onSubmitAddHost = () => {
  console.log('submit add host!')
  AddHost(addHostForm.groupName, addHostForm.ip, addHostForm.hostName).then(result => {
    if (result !== '') {
      ElMessage.error('save failed!' + result)
    } else {
      addHostForm.groupName = ''
      addHostForm.ip = ''
      addHostForm.hostName = ''
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('save successfully!')
    }
  })
}

const onSubmitAllHosts = () => {
  console.log('submit all hosts!')
  SaveAllHosts(allHostsForm.text).then(result => {
    if (result !== '') {
      ElMessage.error('save failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('save successfully!')
    }
  })
}

const onSubmitAllInUseHosts = () => {
  console.log('submit all in use hosts!')
  SaveAllInUseHosts(allInUseHostsForm.text).then(result => {
    if (result !== '') {
      ElMessage.error('save failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('save successfully!')
    }
  })
}

const onSubmitAllGroupHosts = () => {
  console.log('submit all group hosts!')
  SaveAllGroupHosts(groupName.value, allGroupHostsForm.text).then(result => {
    if (result !== '') {
      ElMessage.error('save failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('save successfully!')
    }
  })
}

const onSubmitAddHostsTextForm = () => {
  console.log('submit add host text!')
  SaveAddHostsText(addHostsTextForm.text).then(result => {
    if (result !== '') {
      ElMessage.error('save failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      addHostsTextForm.text = ""
      ElMessage.success('save successfully!')
    }
  })
}

const copyToClipboard = (form) => {
  console.log('copy to clipboard!')
  ClipboardSetText(form.text).then(result => {
    if (!result) {
      ElMessage.error('copy to clipboard failed!')
    } else {
      ElMessage.success('copy to clipboard successfully!')
    }
  })
}

const handleClickTab = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event)
  if (activeMenuIndex.value === 'allHosts') {
    if (tab.props.name === "normal") {
      getAllList()
    } else {
      getHostsText()
    }
  } else if (activeMenuIndex.value === 'addHost') {

  } else if (activeMenuIndex.value === 'inUse') {
    if (tab.props.name === "normal") {
      getActiveList()
    } else {
      getInUseHostsText()
    }
  } else if (activeMenuIndex.value === 'showGroup') {
    getListByGroup()
  }
}

const filterTableData = computed(() =>
    tableData.value.filter(
        (data) => (!filterIP.value || data.ip.toLowerCase().includes(filterIP.value.toLowerCase()))
            && (!filterHostname.value || data.hostName.toLowerCase().includes(filterHostname.value.toLowerCase()))
            && (!filterGroupName.value || data.groupName.toLowerCase().includes(filterGroupName.value.toLowerCase()))
    )
)

const multipleSelection = ref([])
const handleSelectionChange = (val: []) => {
  multipleSelection.value = val
  console.log(multipleSelection)
}

const handleDeleteAllSelected = () => {
  console.log('delete all selected')
  ElMessageBox.confirm('Are you sure to delete all selected?')
      .then(() => {
        DeleteHosts(multipleSelection.value).then(result => {
          if (result !== '') {
            ElMessage.error('delete failed!' + result)
          } else {
            handleMenuSelect(activeMenuIndex.value, null)
            ElMessage.success('delete successfully!')
          }
        })
      })
      .catch(() => {
        // catch error
      })
}

const getAllGroupNames = () => {
  console.log('get all group names')

  GetAllGroupNames().then(result => {
    console.log("group names", result)
    allGroupNames.value = []
    for (let k in result) {
      allGroupNames.value.push({value: result[k]})
    }
  })
}

const querySearchGroupNames = (queryString: string, cb: any) => {
  const results = queryString
      ? allGroupNames.value.filter(
          (groupName) => groupName.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
      )
      : allGroupNames.value
  cb(results)
}

const resetDialogForm = () => {
  dialogForm.index = 0
  dialogForm.id = 0
  dialogForm.oldGroupName = ''
  dialogForm.groupName = ''
}

const openDialogForm = (groupName: string) => {
  dialogFormVisible.value = true
  dialogForm.oldGroupName = groupName
  dialogForm.groupName = groupName
}

const submitDialogForm = () => {
  dialogFormVisible.value = false
  SetGroupName(dialogForm.oldGroupName, dialogForm.groupName).then(result => {
    if (result !== '') {
      ElMessage.error('save failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('save successfully!')
    }
  })
  resetDialogForm()
}

const openDialogMoveForm = (index, group) => {
  dialogMoveFormVisible.value = true
  dialogForm.index = index
  dialogForm.id = group.id
  dialogForm.groupName = group.groupName
  console.log(dialogForm)
}

const submitDialogMoveForm = () => {
  dialogMoveFormVisible.value = false
  SetGroupNameByHostnameId(dialogForm.id, dialogForm.groupName).then(result => {
    if (result !== '') {
      ElMessage.error('save failed!' + result)
    } else {
      tableData.value.splice(dialogForm.index, 1)
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('save successfully!')
    }
  })
  resetDialogForm()
}

const submitDialogMoveSelectedForm = () => {
  console.log('move all selected')
  dialogMoveSelectedFormVisible.value = false
  SetGroupNameByHosts(multipleSelection.value, dialogForm.groupName).then(result => {
    if (result !== '') {
      ElMessage.error('delete failed!' + result)
    } else {
      handleMenuSelect(activeMenuIndex.value, null)
      ElMessage.success('delete successfully!')
    }
  })
  resetDialogForm()
}

onMounted(() => {
  getListByGroup()
  getAllGroupNames()
})
</script>

<style scoped>
.show-content {
  height: 100vh;
}

.ml-2 {
  margin-left: 20px;
}

.mr-1 {
  margin-right: 10px;
}

.mr-2 {
  margin-right: 20px;
}

.mt-2 {
  margin-top: 20px;
}

.border-right {
  border-right: 1px solid var(--el-border-color);
}

.no-border-right {
  border: none;
}
</style>