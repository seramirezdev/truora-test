<template>
    <div id="app" class="h-100">

        <Navbar/>

        <b-container fluid class="my-3">
            <router-view/>
        </b-container>

        <Footer/>

    </div>
</template>

<script>
    import Navbar from "./components/Navbar";
    import Footer from "./components/Footer";

    export default {
        name: 'App',
        components: {Footer, Navbar}
    }
</script>