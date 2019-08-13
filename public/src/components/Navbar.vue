<template>
    <b-navbar toggleable="lg" type="dark" style="background-color: #6666e5; height: 5rem;" class="shadow">
        <b-navbar-brand :to="{ name: 'home' }">Truora Test</b-navbar-brand>

        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse id="nav-collapse" is-nav>
            <b-navbar-nav>
                <b-nav-item :to="{ name: 'home' }" exact>Consultar dominio</b-nav-item>
                <b-nav-item :to="{ name: 'listDomains' }">Lista de dominios</b-nav-item>
            </b-navbar-nav>
        </b-collapse>
    </b-navbar>

</template>

<script>
    export default {
        name: "Navbar"
    }
</script>