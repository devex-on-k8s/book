package dev.minsalus.appointments;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.web.reactive.server.WebTestClient;

import java.time.Instant;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class AppointmentsApplicationTests {

    @Autowired
    private AppointmentRepository appointmentRepository;

    @Autowired
    private WebTestClient webTestClient;

    @BeforeEach
    void setup() {
        appointmentRepository.deleteAll();
    }

    @Test
    void welcomeMessageShouldBeReturned() {
        webTestClient.get()
            .uri("/")
            .exchange()
            .expectStatus().isOk()
            .expectBody(String.class)
            .isEqualTo("Welcome to the Appointments API!");
    }

    @Test
    void appointmentsShouldBeReturned() {
        var appointment = Appointment.with(1L, Instant.now());
        appointmentRepository.save(appointment);

        webTestClient.get()
            .uri("/appointments")
            .exchange()
            .expectStatus().isOk()
            .expectBodyList(Appointment.class)
            .hasSize(1)
            .value(result -> {
                Appointment actualAppointment = result.getFirst();
                assertThat(actualAppointment.id()).isNotNull();
                assertThat(actualAppointment.patientId()).isEqualTo(1L);
                assertThat(actualAppointment.appointmentDate().toEpochMilli()).isEqualTo(appointment.appointmentDate().toEpochMilli());
            });
    }

    @Test
    void appointmentShouldBeCreated() {
        var appointment = Appointment.with(2L, Instant.now());

        webTestClient.post()
            .uri("/appointments")
            .bodyValue(appointment)
            .exchange()
            .expectStatus().isCreated()
            .expectHeader().valueMatches("Location", ".*/appointments/\\d+");

        assertThat(appointmentRepository.findAll()).hasSize(1);
    }

    @Test
	void appointmentsShouldBeDeleted() {
		var appointment = Appointment.with(3L, Instant.now());
		appointmentRepository.save(appointment);

		webTestClient.delete()
            .uri("/appointments")
            .exchange()
            .expectStatus().isNoContent();

		assertThat(appointmentRepository.findAll()).isEmpty();
  }

    @Test
    void appointmentWithCategoryShouldBeCreated() {
        var appointment = new Appointment(null, 4L, "Cardiology", Instant.now());

        webTestClient.post()
            .uri("/appointments")
            .bodyValue(appointment)
            .exchange()
            .expectStatus().isCreated()
            .expectBody(Appointment.class);
    }

}
