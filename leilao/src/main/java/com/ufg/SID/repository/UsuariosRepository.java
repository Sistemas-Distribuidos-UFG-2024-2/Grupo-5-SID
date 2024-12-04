package com.ufg.SID.repository;

import com.ufg.SID.model.Usuarios;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface UsuariosRepository extends JpaRepository<Usuarios, UUID> {
    Optional<Usuarios> findByMail(String mail); // Busca pelo email
}
