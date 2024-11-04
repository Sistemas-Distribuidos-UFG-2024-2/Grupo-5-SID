package com.ufg.SID.repository;

import com.ufg.SID.model.Leilao;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface LeilaoRepository extends JpaRepository<Leilao, Long> {
    List<Leilao> findByParticipantesContains(Long usuarioId);
}
