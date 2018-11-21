<template>
    <v-app id="inspire">
        <v-content>
            <v-container fluid fill-height>
                <v-layout align-center justify-center>
                    <v-flex xs12 sm8 md4>
                        <v-card class="elevation-12">
                            <v-toolbar dark color="primary">
                                <v-toolbar-title>Login</v-toolbar-title>
                                <v-spacer></v-spacer>
                            </v-toolbar>
                            <v-card-text>
                                <v-form>
                                    <p v-if="errors.length">
                                        <v-alert v-for="error in errors"
                                                 :value="true"
                                                 type="error">
                                            {{ error }}
                                        </v-alert>
                                    </p>

                                    <v-text-field v-model="input.username" prepend-icon="person" name="login"
                                                  label="Login"
                                                  type="text"></v-text-field>
                                    <v-text-field v-model="input.password" prepend-icon="lock" name="password"
                                                  label="Password" id="password"
                                                  type="password"></v-text-field>
                                </v-form>
                            </v-card-text>
                            <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn color="primary" v-on:click.native="login()">Login</v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-container>
        </v-content>
    </v-app>
</template>

<script>

    export default {

        name: 'Login',
        data() {
            return {
                drawer: null,
                errors: [],
                input: {
                    username: "",
                    password: ""
                }
            }
        },
        methods: {
            login() {
                if (!this.input.username) {
                    this.errors.push('Username required.');
                }
                if (!this.input.password) {
                    this.errors.push('Password required.');
                }

                if (this.input.username && this.input.password) {
                    if (this.input.username === this.$parent.mockAccount.username && this.input.password === this.$parent.mockAccount.password) {
                        this.$emit("authenticated", true);
                        this.$router.replace({name: "secure"});
                    } else {
                        this.errors.push("The username and / or password is incorrect");
                    }
                }


            }
        },
        props: {
            source: String
        }
    }
</script>
