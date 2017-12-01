let store = new Vuex.Store({
    state: {
        taskList: []
    },
    mutations: {
        loadTasks(state) {
            fetch('/tasks')
                .then((response) => {
                    if(response.ok) {
                        return response.json();
                    }
                    throw new Error('Network response was not ok');
                })
                .then((json) => {
                    if (json.items){
                        json.items.forEach(function(task) {
                            store.commit('addTask', task);
                        });
                    }
                })
                .catch((error) => {
                    console.log(error);
                });
        },
        addTask (state, task) {
            state.taskList.push(task);

        },
        createTask (state, task) {
            //let json = new FormData();
            //json.append( "json", JSON.stringify( task ) );
            let json = JSON.stringify( task );
            console.log(json);

            fetch('/tasks', {
                method: 'PUT',
                headers: {
                    'Accept': 'application/json, text/plain, */*',
                    'Content-Type': 'application/json'
                },
                body: json
            })
                .then(function(response){
                    if(response.ok) {
                        return response.json();
                    }
                    throw new Error('Network response was not ok');
                })
                .then(function(json){
                    task.id = json.created;
                    store.commit('addTask', task);
                    console.log("Task created!")
                })
                .catch(function(error){
                    console.log(error);
                });
        },
        delTask (state, task) {
            fetch('/tasks/' + task.id, {
                method: 'DELETE'
            })
                    .then((response) => {
                        if(response.ok) {
                            return response.json();
                        }
                        throw new Error('Network response was not ok');
                    })
                    .then((json) => {
                        const taskIndex = state.taskList.indexOf(task);
                        state.taskList.splice(taskIndex, 1);
                        console.log("Task deleted!");
                    })
                    .catch((error) => {
                        console.log(error);
                    });
        }
    }
});


document.addEventListener("DOMContentLoaded", function(){
    new Vue({
        el: '#app',
        data: {
            taskList: store.state.taskList
        },
        created: function() {
            store.commit('loadTasks');
        }
    })
});

Vue.component('task-item', {
    props: ['task'],
    template:
    `<li class="list-group-item">
        <span>{{ task.name }}</span>
        <button class="btn btn-del" @click="deleteTask(task)"><i class="fa fa-trash-o"></i></button>
    </li>`,
    methods: {
        deleteTask(task) {
            store.commit('delTask', task);
        }
    }
});

Vue.component('add-task', {
    template:
    `<div class="el-card">
        <div class="input-group">
            <input type="text" placeholder="New Task" class="form-control" v-model="task.name"/>
            <span class="input-group-btn">
                <button class="btn btn-create" v-on:click="createTask(task)">Create</button>
            </span>
        </div>
    </div>`,
    data() {
        return {
            task: {
                id: 0,
                name: ""
            }
        }
    },
    methods: {
        createTask(task) {
            store.commit('createTask', task);
            this.task = {};
        }
    }
});