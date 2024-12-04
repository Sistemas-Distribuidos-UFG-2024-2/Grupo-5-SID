package com.ufg.SID.controller;

import com.ufg.SID.model.Usuarios;
import com.ufg.SID.repository.UsuariosRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Optional;

@RestController
@RequestMapping("/usuarios")
public class UsuariosController {

    @Autowired
    private UsuariosRepository usuariosRepository;

    @PostMapping("/login")
    public ResponseEntity<?> login(@RequestBody Usuarios usuario) {
        Optional<Usuarios> user = usuariosRepository.findByMail(usuario.getMail());
        if (user.isPresent()) {
            return ResponseEntity.ok().body("{\"id\": \"" + user.get().getId().toString() + "\", \"mail\": \"" + user.get().getMail() + "\"}");

        } else {
            return ResponseEntity.status(404).body("{\"error\": \"Usuário não encontrado\"}");
        }
    }
}
