<template>
    <div>

        <b-card v-if="domains.length === 0" class="shadow-sm text-center" bg-variant="danger" text-variant="white">
            <h1>No se han consultados dominios</h1>
        </b-card>
        <b-row v-else>
            <b-col cols="12" md="12" lg="6" v-for="(domain, index) of domains" class="mb-2" :key="index">
                <b-card no-body class="mb-0 p-0 shadow">
                    <b-card-header header-tag="header" class="p-0" role="tab">
                        <b-button block href="#" v-b-toggle="'accordion-'+index" variant="info">
                            <h3><img :src="domain.logo" alt="" height="30"> {{ domain.name }}</h3>
                        </b-button>
                    </b-card-header>
                    <b-collapse :id="'accordion-'+index" accordion="my-accordion" role="tabpanel" visible>
                        <b-card-body>
                            <b-card-text>
                                <p>
                                    <b>{{ domain.title }}</b>
                                </p>
                                <p>
                                    <b :class="domain.is_down ? 'text-danger' : 'text-success'">
                                        {{ domain.is_down ? 'Servidores caidos' : 'Servidores activos' }}
                                    </b>
                                </p>
                                <p>
                                    <b>Grado SSL más bajo: </b>{{ domain.ssl_grade }}
                                </p>
                                <p>
                                    <b>Grado SSL previo: </b>{{ domain.ssl_grade }}
                                </p>
                                <p>
                                    <b>Servidores: </b>
                                    <b-badge pill variant="danger">{{ domain.servers.length }}</b-badge>
                                </p>

                                <SeversInfo v-bind:servers="domain.servers"/>

                            </b-card-text>

                        </b-card-body>
                    </b-collapse>
                </b-card>
            </b-col>
        </b-row>

    </div>
</template>

<script>

    import axios from 'axios';
    import SeversInfo from "../components/SeversInfo";

    export default {
        name: "DomainsList",
        components: {SeversInfo},
        data() {
            return {
                domains: []
            };
        },
        async mounted() {

            try {
                let json = await axios.get('http://localhost:3000/domains');

                this.domains = json.data.items;
            } catch (e) {
                console.log(e);
            }
        }
    }

</script>