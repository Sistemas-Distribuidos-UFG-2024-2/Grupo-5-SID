package com.ufg.SID.model;
import jakarta.persistence.*;

@Entity
public class Usuarios {

        @Id
        @GeneratedValue(strategy = GenerationType.IDENTITY)
        private Long cod;

        private String nome;
        private String email;
        private String senha;


        public String getSenha() {
                return senha;
        }
        public void setSenha(String senha) {
                this.senha = senha;
        }
        public String getEmail() {
                return email;
        }
        public void setEmail(String email) {
                this.email = email;
        }
        public Long getCod() {
                return cod;
        }
        public void setCod(Long cod) {
                this.cod = cod;
        }
        public String getNome() {
                return nome;
        }
        public void setNome(String nome) {
                this.nome = nome;
        }
        

}
