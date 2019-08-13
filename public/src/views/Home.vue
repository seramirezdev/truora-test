<template>
    <div>
        <div class="shadow">
            <b-card-header header-bg-variant="info" header-text-variant="white">
                <h3>Obtén información de diferentes dominio!</h3>
            </b-card-header>
            <b-card-body>
                <b-row>
                    <b-col cols="10">
                        <label class="sr-only">Name</label>
                        <b-form-group
                                description="Esta consulta puede tardar algunos segundos, obtenemos la información de diferentes sitios">
                            <b-input-group prepend="www." class="mb-2">
                                <b-input
                                        v-model.string="domainToSearch"
                                        @keyup.enter="getDomain"
                                        placeholder="Ex: truora.com, google.com, amazon.com"
                                        des
                                ></b-input>
                            </b-input-group>
                            <small v-if="invalidDomain" class="text-danger">{{ invalidDomainMsg }}</small>
                        </b-form-group>
                    </b-col>
                    <b-col>
                        <b-button class="btn-block" variant="info" @click="getDomain" :disabled="isLoading">Consultar
                        </b-button>
                    </b-col>

                </b-row>
            </b-card-body>
        </div>
        <b-row>
            <b-col class="mt-3">
                <b-progress v-if="isLoading" :value="100" striped variant="warning" animated/>
            </b-col>
        </b-row>

        <b-row>
            <DomainInfo/>
        </b-row>
    </div>
</template>

<script>
    // @ is an alias to /src
    import axios from 'axios'
    import DomainInfo from "../components/DomainInfo";
    import {mapMutations} from 'vuex';
    import SeverInfo from "../components/SeversInfo";

    export default {
        name: 'home',
        components: {SeverInfo, DomainInfo},
        data() {
            return {
                isLoading: false,
                domainToSearch: '',
                invalidDomain: false,
                invalidDomainMsg: ''
            }
        },
        methods: {
            ...mapMutations(['domainFound']),
            async getDomain() {

                this.domainFound(null);
                if (this.domainToSearch === '') {
                    this.invalidDomain = true;
                    this.invalidDomainMsg = 'Debes ingresar un dominio';
                    return
                }

                this.invalidDomain = false;
                this.invalidDomainMsg = '';
                this.isLoading = true;

                let time = 500;
                let allServersGetSSL = true;


                try {
                    let getAllData = false;
                    const search = this.domainToSearch;
                    do {
                        let json = await axios.get(`http://localhost:3000/consult-domain/${search}`);
                        let domain = json.data;

                        domain = {
                            ...domain,
                            name: this.domainToSearch,
                        };

                        if (domain.servers.length > 0) {

                            for (var s of domain.servers) {
                                if (s.ssl_grade === '') {
                                    allServersGetSSL = false;
                                    break;
                                } else {
                                    allServersGetSSL = true;
                                }

                            }
                            if (allServersGetSSL) {

                                this.isLoading = false;
                                getAllData = true;
                            }
                            time = 1000;
                        }
                        this.domainFound(domain);
                        await sleep(time);
                    } while (!getAllData);
                } catch (e) {
                    console.log(e);
                    this.invalidDomain = true;
                    this.invalidDomainMsg = 'No se encontró el dominio';
                } finally {
                    this.isLoading = false;
                }
            }
        },
    }

    function sleep(time) {
        return new Promise(resolve => setTimeout(resolve, time));
    }
</script>
