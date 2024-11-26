package com.ufg.SID.repository;

import com.ufg.SID.model.Leilao;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface LeilaoRepository extends JpaRepository<Leilao, Long> {
    @Query("SELECT l FROM Leilao l JOIN l.participantes p WHERE p.usuarioEmail = :usuarioEmail")
    List<Leilao> findByUsuarioParticipante(@Param("usuarioEmail") String usuarioEmail);

}
