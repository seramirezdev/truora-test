<template>
    <b-container fluid class="mb-5">
        <div v-if="domainData != null" class="shadow my-3">

            <div v-if="!domainData.invalidDomain">
                <b-card-header
                        class="m-0"
                        header-class="text-center"
                        :header-bg-variant="colorHeader"
                        header-text-variant="white">
                    <h4 v-if="msgHeader !== ''">{{ msgHeader }}</h4>
                </b-card-header>
                <b-card-body>
                    <h3 class="mb-3">{{ domainData.title }}</h3>
                    <img :src="domainData.logo" alt="" height="40"/>

                    <h3 class="ml-3" style="display: inline-flex">{{ domainData.name }}</h3>
                    <h3 :class="domainData.is_down ? 'text-warning': 'text-success'">
                        {{ domainData.is_down ? 'Servidores caido' : 'Servidores funcionando' }}
                    </h3>
                    <h3>
                        Cambios en el dominio: {{ domainData.servers_changed ? 'Sí' : 'No' }}
                    </h3>
                    <h3>
                        Grado SSL más bajo: {{ domainData.ssl_grade }}
                    </h3>
                    <h3>
                        Grado SSL previo: {{ domainData.previous_ssl_grade }}
                    </h3>

                    <h3>
                        Servidores:
                        <b-badge pill :variant="domainData.servers.length === 0 ? 'warning' : 'info   '">
                            {{ domainData.servers.length > 0 ? domainData.servers.length :
                            'Buscando...'}}
                        </b-badge>
                    </h3>

                    <SeversInfo v-bind:servers="domainData.servers"/>

                </b-card-body>
            </div>
        </div>
    </b-container>
</template>

<script>

    import {mapState} from 'vuex';
    import SeversInfo from "./SeversInfo";

    export default {
        name: "DomainInfo",
        components: {SeversInfo},
        data() {
            return {
                serversFound: false,
                sslGradesObtained: false,
                colorHeader: 'warning',
                msgHeader: null,
            }
        },
        computed: {
            ...mapState({
                domainData(state) {
                    if (state.domainData === null) {
                        this.serversFound = false;
                        this.sslGradesObtained = false;
                        this.msgHeader = null;
                        return;
                    }

                    let domain = state.domainData;

                    if (domain.servers.length > 0) {
                        this.serversFound = true;


                        for (var s of domain.servers) {
                            if (s.ssl_grade === '') {
                                this.sslGradesObtained = false;
                                break;
                            } else {
                                this.sslGradesObtained = true;
                            }
                        }
                    }

                    if (!this.serversFound) {
                        this.msgHeader = 'Buscando los servidores del dominio ' + domain.name + ' ...';
                    } else if (!this.sslGradesObtained) {
                        this.msgHeader = 'Obteniendo grado SSL de servidores ' + domain.name + ' ...';
                    } else {
                        this.msgHeader = 'Se obtuvieron todos los datos!';
                    }

                    this.colorHeader = this.serversFound && this.sslGradesObtained ? 'success' : 'warning';

                    return state.domainData;
                }
            })
        }
    }
</script>
