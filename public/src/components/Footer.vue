<template>
    <b-navbar id="footer" type="dark" style="background-color: #6666e5;">
        <b-container fluid>

            <b-nav-text style="font-size: large">&copy; Sergio Ramírez 2019 - seramirezdev@gmail.com</b-nav-text>

        </b-container>
    </b-navbar>
</template>

<script>

    import {mapState} from 'vuex'

    export default {
        name: "Footer",
        computed: {
            ...mapState(['domainData']
            )
        }
    }
</script>
<style scoped>
    #footer {
        position: fixed;
        left:0;
        bottom: 0;
        width: 100%;
    }
</style>