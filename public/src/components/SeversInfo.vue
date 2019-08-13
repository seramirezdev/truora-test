<template>
    <b-row>
        <b-col cols="12" class="mb-3" md="12" lg="6" xl="4" v-for="domain of domains" :key="domains.address">
            <b-card border-variant="info" header-bg-variant="info" :header="'IP: ' + domain.address"
                    header-text-variant="white"
                    header-class="text-center">
                <b-card-body>

                    <b>Grado SSL: </b>
                    <b-badge pill :variant="domain.ssl_grade === '' ? 'warning' : 'info'">{{ domain.ssl_grade === ''
                        ? 'Buscando...' : domain.ssl_grade }}
                    </b-badge>
                    <br>
                    <b>Propietario: </b> {{ domain.owner }}<br>
                    <b>Pais:</b> {{ domain.country }}

                </b-card-body>
            </b-card>
        </b-col>

    </b-row>
</template>

<script>
    export default {
        name: "SeversInfo",
        props: {
            domains: Array
        }
    }
</script>