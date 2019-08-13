<template>
    <b-row>
        <b-col cols="12" class="mb-3" md="12" lg="6" xl="4" v-for="server of servers" :key="server.address">
            <b-card border-variant="info" header-bg-variant="info" :header="'IP: ' + server.address" header-tag="b"
                    header-text-variant="white"
                    header-class="text-center">
                <b-card-body>

                    <b>Grado SSL: </b>
                    <b-badge pill :variant="server.ssl_grade === '' ? 'warning' : 'info'">{{ server.ssl_grade === ''
                        ? 'Buscando...' : server.ssl_grade }}
                    </b-badge>
                    <br>
                    <b>Propietario: </b> {{ server.owner }}<br>
                    <b>Pais:</b> {{ server.country }}

                </b-card-body>
            </b-card>
        </b-col>

    </b-row>
</template>

<script>
    export default {
        name: "SeversInfo",
        props: {
            servers: Array
        }
    }
</script>