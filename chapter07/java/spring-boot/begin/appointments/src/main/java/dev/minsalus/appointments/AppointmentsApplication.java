package dev.minsalus.appointments;

import java.time.Instant;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.tool.annotation.Tool;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.data.annotation.Id;
import com.fasterxml.jackson.annotation.JsonIgnore;

import org.springframework.data.repository.ListCrudRepository;
import org.springframework.stereotype.Component;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.function.RouterFunction;
import org.springframework.web.servlet.function.RouterFunctions;
import org.springframework.web.servlet.function.ServerResponse;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;


@SpringBootApplication
public class AppointmentsApplication {

	public static void main(String[] args) {
		SpringApplication.run(AppointmentsApplication.class, args);
	}

    @Bean
    RouterFunction<ServerResponse> routes(AppointmentRepository appointmentRepository) {
        return RouterFunctions.route()
            .GET("/", _ -> ServerResponse.ok().body("Welcome to the Appointments API!"))
            .GET("/appointments", _ -> 
                ServerResponse.ok().body(appointmentRepository.findAll()))
            .POST("/appointments", request -> {
                var appointment = appointmentRepository.save(request.body(Appointment.class));
                var location = request.uriBuilder().path("/{id}").build(appointment.id());
                return ServerResponse.created(location).body(appointment);
            })
            .DELETE("/appointments", _ -> {
                appointmentRepository.deleteAll();
                return ServerResponse.noContent().build();
            })
            .build();
    }

}

record Appointment(@Id @JsonIgnore Long id, Long patientId, String category, Instant appointmentDate) {
    public static Appointment with(Long patientId, Instant appointmentDate) {
        return new Appointment(null, patientId, "General", appointmentDate);
    }
}

interface AppointmentRepository extends ListCrudRepository<Appointment, Long> {
    List<Appointment> findByCategory(String category);
}

@RestController
class ChatController {

    private final ChatClient chatClient;

    ChatController(ChatClient.Builder chatClientBuilder, AppointmentTools appointmentTools) {
        this.chatClient = chatClientBuilder
            .defaultTools(appointmentTools)
            .build();
    }

    @PostMapping("chat")
    String postMethodName(@RequestBody String question) {
        return chatClient.prompt(question)
            .call()
            .content();
    }

}   

@Component
class AppointmentTools {

    private static final Logger logger = LoggerFactory.getLogger(AppointmentTools.class);
    private final AppointmentRepository appointmentRepository;

    AppointmentTools(AppointmentRepository appointmentRepository) {
        this.appointmentRepository = appointmentRepository;
    }

    @Tool(description = "List all scheduled appointments of a certain category/type")
    List<Appointment> listAppointmentsByCategory(String category) {
        logger.info("Listing all {} appointments", category);
        return appointmentRepository.findAll();
    }

}
